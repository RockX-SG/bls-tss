#![feature(proc_macro_hygiene)]
use bls_eth::threshold_bls::party_i::LocalKey;
use bls_eth::threshold_bls::state_machine::keygen::{Keygen, R};
use bls_eth::threshold_bls::state_machine::sign::{Sign};
use round_based::{Msg, StateMachine};
use std::convert::From;
use std::ffi::CStr;
use std::fmt::{Debug, Display};
use concat_idents::concat_idents;
use anyhow::Result;

const STATUS_OK: i32 = -0x0000;
const ERROR_STATE_IS_NULL: i32 = -0x1001;
const ERROR_NULL_OR_EMPTY_VALUE: i32 = -0x2001;
const ERROR_STATE_MACHINE_INTERNAL_ERROR: i32 = -0x3001;
const ERROR_INTEROP_BUFFER_TOO_SMALL_ERROR: i32 = -0x4001;
const ERROR_MESSAGE_SERDE_ERROR: i32 = -0x5001;

trait ToI32: Sized {
    fn to_i32(self) -> i32;
}

impl ToI32 for Option<u16> {
    fn to_i32(self) -> i32 {
        match self {
            Some(num) => {
                i32::from(num)
            }
            None => { ERROR_NULL_OR_EMPTY_VALUE }
        }
    }
}

impl ToI32 for bool {
    fn to_i32(self) -> i32 {
        i32::from(self)
    }
}

impl ToI32 for u16 {
    fn to_i32(self) -> i32 {
        i32::from(self)
    }
}

trait StateMachineOutput {
    fn pick_string_output(&mut self) -> (Option<String>, i32);
}

trait StateMachineState {
    fn get_string_state(&mut self) -> (Option<String>, i32);
}

impl StateMachineState for Keygen {
    fn get_string_state(&mut self) -> (Option<String>, i32) {
        let machine_state = self.get_state();
        let res = serde_json::to_string(&machine_state).unwrap_or_default();
        (Some(res), STATUS_OK)
    }
}

impl StateMachineOutput for Keygen {
    fn pick_string_output(&mut self) -> (Option<String>, i32) {
        match self.pick_output() {
            Some(Ok(res)) => {
                let res = serde_json::to_string(&res).unwrap_or_default();
                (Some(res), STATUS_OK)
            }
            Some(Err(_)) => { (None, ERROR_STATE_MACHINE_INTERNAL_ERROR) }
            None => { (None, STATUS_OK) }
        }
    }
}

impl StateMachineOutput for Sign {
    fn pick_string_output(&mut self) -> (Option<String>, i32) {
        match self.pick_output() {
            Some(Ok((_, sig))) => {
                let compressed = true;
                let bytes = sig.to_bytes(compressed);
                let sig = hex::encode(bytes.as_slice());
                (Some(sig), STATUS_OK)
            }
            Some(Err(_)) => { (None, ERROR_STATE_MACHINE_INTERNAL_ERROR) }
            None => { (None, STATUS_OK) }
        }
    }
}

fn write_to_buffer(output: &String, buf: *mut cty::c_char, max_len: cty::c_int) -> cty::c_int {
    let src = output.as_bytes().as_ptr();
    let len = output.as_bytes().len();
    let len_c_int = len as cty::c_int;
    if len_c_int <= max_len - 1 {
        unsafe {
            std::ptr::copy(src, buf as *mut u8, len);
            (*buf.offset(len as isize)) = 0;
        }
        len_c_int
    } else {
        ERROR_INTEROP_BUFFER_TOO_SMALL_ERROR
    }
}

fn ret_or_err<T, E>(res: Result<T, E>) -> *mut T where E: Debug + Display {
    match res {
        Ok(res) => { Box::into_raw(Box::new(res)) }
        Err(e) => {
            log::error!("Encountered error: {}", e);
            std::ptr::null_mut()
        }
    }
}

macro_rules! create_function {
    // This macro takes an argument of designator `ident` and
    // creates a function named `$func_name`.
    // The `ident` designator is used for variable/function names.
    ($sm_type:ty,$sm_name:ident,$func_name:ident) => {
        concat_idents!(full_name=$sm_name, _, $func_name, {
            #[no_mangle]
            pub extern "C" fn full_name(state: Option<&$sm_type>) -> cty::c_int {
                match state {
                    Some(state) => { state.$func_name().to_i32() }
                    None => { ERROR_STATE_IS_NULL }
                }
            }
        });
    };
}

macro_rules! create_free_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=free, _, $sm_name, {
            #[no_mangle]
            pub unsafe extern "C" fn full_name(state: *mut $sm_type) {
                assert!(!state.is_null());
                Box::from_raw(state); // Rust auto-drops it
            }
        });
    };
}

macro_rules! create_has_outgoing_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=$sm_name, _, has_outgoing, {
            #[no_mangle]
            pub extern "C" fn full_name(state: Option<& mut $sm_type>) -> cty::c_int {
                match state {
                    Some(state) => { state.message_queue().len() as cty::c_int }
                    None => { ERROR_STATE_IS_NULL }
                }
            }
        });
    };
}

macro_rules! create_proceed_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=$sm_name, _, proceed, {
            #[no_mangle]
            pub unsafe extern "C" fn full_name(state: Option<&mut $sm_type>) -> cty::c_int {
                match state {
                    Some(state) => {
                        match state.proceed() {
                            Ok(_) => {STATUS_OK}
                            Err(e) => {
                                log::error!("Failed to proceed: {}", e);
                                ERROR_STATE_MACHINE_INTERNAL_ERROR
                            }
                        }

                    }
                    None => { ERROR_STATE_IS_NULL }
                }
            }
        });
    };
}

macro_rules! create_incoming_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=$sm_name, _, incoming, {
            #[no_mangle]
            pub extern "C" fn full_name(state: Option<&mut $sm_type>, buf: *const cty::c_char) -> cty::c_int {
                match state {
                    Some(state) => {
                        let arr = unsafe { CStr::from_ptr(buf).to_bytes() };
                        let res = serde_json::from_slice::<Msg<<$sm_type as StateMachine>::MessageBody>>(arr);
                        match res {
                            Ok(msg) => {
                                let h_res = state.handle_incoming(msg);
                                match h_res {
                                    Ok(_) => {
                                        STATUS_OK
                                    }
                                    Err(e) => {
                                        log::error!("Failed to handle incoming message: {}", e);
                                        ERROR_STATE_MACHINE_INTERNAL_ERROR
                                    }
                                }
                            }
                            Err(e) => {
                                log::error!("Failed to parse incoming message: {}", e);
                                ERROR_MESSAGE_SERDE_ERROR
                            }
                        }
                    }
                    None => {
                        ERROR_STATE_IS_NULL
                    }
                }
            }
        });
    };
}

macro_rules! create_outgoing_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=$sm_name, _, outgoing, {

            #[no_mangle]
            pub unsafe extern "C" fn full_name(state: Option<&mut $sm_type>, buf: *mut cty::c_char, max_len: cty::c_int) -> cty::c_int {
                match state {
                    Some(state) => {
                        let msg = state.message_queue().drain(..1).next();
                        match msg {
                            Some(msg) => {
                                let res = serde_json::to_string(&msg);
                                match res {
                                    Ok(str) => {
                                        write_to_buffer(&str, buf, max_len)
                                    }
                                    Err(_) => {
                                        ERROR_MESSAGE_SERDE_ERROR
                                    }
                                }
                            }
                            None => { STATUS_OK }
                        }
                    }
                    None => { ERROR_STATE_IS_NULL }
                }
            }
        });
    };
}

macro_rules! create_pick_output_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=$sm_name, _, pick_output, {

            #[no_mangle]
            pub unsafe extern "C" fn full_name(state: Option<&mut $sm_type>, buf: *mut cty::c_char, max_len: cty::c_int) -> cty::c_int {
                match state {
                    Some(state) => {
                        let output = state.pick_string_output();
                        match output {
                            (Some(str), _) => {
                                write_to_buffer(&str, buf, max_len)
                            }
                            (None, status) => {
                                status
                            }
                        }
                    }
                    None => { ERROR_STATE_IS_NULL }
                }
            }
        });
    };
}

macro_rules! create_get_state_function {
    ($sm_type:ty,$sm_name:ident) => {
        concat_idents!(full_name=$sm_name, _, get_state, {

            #[no_mangle]
            pub unsafe extern "C" fn full_name(state: Option<&mut $sm_type>, buf: *mut cty::c_char, max_len: cty::c_int) -> cty::c_int {
                match state {
                    Some(state) => {
                        let machine_state = state.get_string_state();
                        match machine_state {
                            (Some(str), _) => {
                                write_to_buffer(&str, buf, max_len)
                            }
                            (None, status) => {
                                status
                            }
                        }
                    }
                    None => { ERROR_STATE_IS_NULL }
                }
            }
        });
    };
}

macro_rules! create_wrapper {
    ($sm_type:ty,$sm_name:ident) => {

        create_free_function!($sm_type, $sm_name);

        create_has_outgoing_function!($sm_type, $sm_name);

        create_proceed_function!($sm_type, $sm_name);

        create_incoming_function!($sm_type, $sm_name);

        create_outgoing_function!($sm_type, $sm_name);

        create_pick_output_function!($sm_type, $sm_name);

        create_function!($sm_type, $sm_name, total_rounds);

        create_function!($sm_type, $sm_name, current_round);

        create_function!($sm_type, $sm_name, party_ind);

        create_function!($sm_type, $sm_name, parties);

        create_function!($sm_type, $sm_name, is_finished);

        create_function!($sm_type, $sm_name, wants_to_proceed);

    };
}

create_wrapper!(Keygen, keygen);

create_wrapper!(Sign, sign);

create_get_state_function!(Keygen, keygen);

#[no_mangle]
pub extern "C" fn new_keygen(i: cty::c_int, t: cty::c_int, n: cty::c_int) -> *mut Keygen {
    let state = Keygen::new(i as u16, t as u16, n as u16);
    ret_or_err(state)
}

#[no_mangle]
pub extern "C" fn keygen_from_state(round: *const cty::c_char) -> *mut Keygen {
    let round = unsafe { CStr::from_ptr(round).to_bytes() };
    let round = serde_json::from_slice::<R>(round);
    match round {
        Ok(round) => {
            let state = Keygen::from(round);
            let state: Result<Keygen> = Ok(state);
            ret_or_err(state)
        }
        Err(e) => {
            log::error!("Failed to decode the round: {}", e);
            std::ptr::null_mut()
        }
    }
}

#[no_mangle]
pub extern "C" fn new_sign(message_hash: *const cty::c_char, i: cty::c_int, n: cty::c_int, local_key: *const cty::c_char) -> *mut Sign {
    let message_hash = unsafe { CStr::from_ptr(message_hash).to_bytes() };
    let message_hash = message_hash.to_vec();
    let local_key = unsafe { CStr::from_ptr(local_key).to_bytes() };
    let local_key = serde_json::from_slice::<LocalKey>(local_key);

    match local_key {
        Ok(local_key) => {
            let state = Sign::new(message_hash, i as u16, n as u16, local_key);
            ret_or_err(state)
        }
        Err(e) => {
            log::error!("Failed to decode the local key: {}", e);
            std::ptr::null_mut()
        }
    }
}
