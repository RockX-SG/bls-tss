#![feature(proc_macro_hygiene)]
use bls::threshold_bls::state_machine::keygen::{Keygen, LocalKey};
use round_based::{Msg, StateMachine};
use std::convert::From;
use std::ffi::CStr;
use std::ffi::CString;
use cty::c_char;
use concat_idents::concat_idents;

type KeygenMsg = Msg<<Keygen as StateMachine>::MessageBody>;

macro_rules! create_function {
    // This macro takes an argument of designator `ident` and
    // creates a function named `$func_name`.
    // The `ident` designator is used for variable/function names.
    ($sm_type:ty,$sm_name:ident,$func_name:ident) => {
        concat_idents!(full_name=$sm_name, _, $func_name, {
            #[no_mangle]
            pub extern "C" fn full_name(state: Option<&$sm_type>) -> cty::c_int {
                match state {
                    Some(state) => { cty::c_int::from(state.$func_name()) }
                    None => { -1 }
                }
            }
        });
    };
}


#[no_mangle]
pub extern "C" fn new_keygen(i: cty::c_int, t: cty::c_int, n: cty::c_int) -> *mut Keygen {
    let state = Keygen::new(i as u16, t as u16, n as u16);
    match state {
        Ok(state) => { Box::into_raw(Box::new(state)) }
        Err(e) => {
            println!("error: {:?}", e);
            std::ptr::null_mut()
        }
    }
}

#[no_mangle]
pub unsafe extern "C" fn free_keygen(state: *mut Keygen) {
    assert!(!state.is_null());
    Box::from_raw(state); // Rust auto-drops it
}

#[no_mangle]
pub extern "C" fn keygen_total_rounds(state: Option<&Keygen>) -> cty::c_int {
    match state {
        Some(state) => {
            match state.total_rounds() {
                Some(tr) => { cty::c_int::from(tr) }
                None => { -1 }
            }
        }
        None => { -1 }
    }
}

create_function!(Keygen, keygen, current_round);

create_function!(Keygen, keygen, party_ind);

create_function!(Keygen, keygen, parties);

create_function!(Keygen, keygen, is_finished);

create_function!(Keygen, keygen, wants_to_proceed);


#[no_mangle]
pub extern "C" fn keygen_has_outgoing(state: Option<& mut Keygen>) -> cty::c_int {
    match state {
        Some(state) => { state.message_queue().len() as cty::c_int }
        None => { -1 }
    }
}

#[no_mangle]
pub unsafe extern "C" fn keygen_proceed(state: Option<&mut Keygen>) -> cty::c_int {
    match state {
        Some(state) => {
            match state.proceed() {
                Ok(_) => {0}
                Err(e) => {
                    println!("error: {:?}", e);
                    -2
                }
            }

        }
        None => { -1 }
    }
}

#[no_mangle]
pub unsafe extern "C" fn keygen_incoming(state: Option<&mut Keygen>, buf: *const cty::c_char) -> cty::c_int {
    match state {
        Some(state) => {
            let arr = unsafe { CStr::from_ptr(buf).to_bytes() };
            let res = serde_json::from_slice::<KeygenMsg>(arr);
            match res {
                Ok(msg) => {
                    let hRes = state.handle_incoming(msg);
                    match hRes {
                        Ok(_) => {
                            0
                        }
                        Err(e) => {
                            println!("error: {:?}", e);
                            -1
                        }
                    }
                }
                Err(e) => {
                    println!("error: {:?}", e);
                    -1
                }
            }
        }
        None => {
            -1
        }
    }
}

#[no_mangle]
pub unsafe extern "C" fn keygen_outgoing(state: Option<&mut Keygen>, buf: *mut cty::c_char, maxlen: cty::c_int) -> cty::c_int {
    match state {
        Some(state) => {
            let msg = state.message_queue().drain(..1).next();
            match msg {
                Some(msg) => {
                    let res = serde_json::to_string(&msg);
                    match res {
                        Ok(str) => {
                            write_to_buffer(&str, buf, maxlen)
                        }
                        Err(e) => {
                            -2
                        }
                    }
                }
                None => { 0 }
            }
        }
        None => { -1 }
    }
}

unsafe fn write_to_buffer(output:&String, buf: *mut cty::c_char, maxlen: cty::c_int) -> cty::c_int {
    let src = output.as_bytes().as_ptr();
    let len = output.as_bytes().len();
    let len_c_int = len as cty::c_int;
    if len_c_int <= maxlen - 1 {
        unsafe {
            std::ptr::copy(src, buf as *mut u8, len);
            (*buf.offset(len as isize)) = 0;
        }
        len_c_int
    } else {
        -3
    }
}

#[no_mangle]
pub unsafe extern "C" fn keygen_pick_output(state: Option<&mut Keygen>, buf: *mut cty::c_char, maxlen: cty::c_int) -> cty::c_int {
    match state {
        Some(state) => {
            let output = state.pick_output();
            match output {
                Some(Ok(localKey)) => {
                    let res = serde_json::to_string(&localKey);
                    match res {
                        Ok(str) => {
                            write_to_buffer(&str, buf, maxlen)
                        }
                        Err(e) => {
                            -2
                        }
                    }
                }
                Some(Err(e)) => {
                    println!("error: {:?}", e);
                    -1
                }
                None => {
                    println!("error: Not finished yet.");
                    -1
                }
            }
        }
        None => { -1 }
    }
}