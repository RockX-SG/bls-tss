package bls_tss

/*
#cgo CFLAGS:-I${SRCDIR}/bls-tss/include
#include <stdio.h>
#include <stdlib.h>
#include <tss.h>
*/
import "C"
import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"unsafe"
)

type KeygenSimple struct {
	i      int
	t      int
	n      int
	state  unsafe.Pointer
	buffer unsafe.Pointer
	output *string
}

func NewKeygenSimple(i int, t int, n int) *KeygenSimple {
	buffer := C.malloc(C.size_t(BufferSize))
	state := C.new_keygen(C.int(i), C.int(t), C.int(n))
	return &KeygenSimple{i, t, n, state, buffer, nil}
}

func NewKeygenSimpleFromRound(i, t, n int, round string) *KeygenSimple {
	buffer := C.malloc(C.size_t(BufferSize))
	cRound := C.CString(round)
	defer C.free(unsafe.Pointer(cRound))
	state := C.keygen_from_state(cRound)
	return &KeygenSimple{i, t, n, state, buffer, nil}
}

func (k *KeygenSimple) Init() []string {
	k.proceedIfNeeded()
	return k.getOutgoing()
}

func (k *KeygenSimple) Handle(msg string) (bool, []string, error) {
	k.handleIncoming(msg)
	k.proceedIfNeeded()
	outgoing := k.getOutgoing()
	output := k.finishIfPossible()
	finished := output != nil

	if finished {
		k.output = output
	}

	return finished, outgoing, nil
}

func (k *KeygenSimple) Free() {
	C.free(k.buffer)
	C.free_keygen(k.state)
}

func (k *KeygenSimple) Output() *string {
	return k.output
}

func (k *KeygenSimple) proceedIfNeeded() {
	res := C.keygen_wants_to_proceed(k.state)
	k.trace("keygen_wants_to_proceed", res)
	if res == 1 {
		res = C.keygen_proceed(k.state)
		k.trace("keygen_proceed", res)
	}
}

func (k *KeygenSimple) getOutgoing() []string {
	var outgoing []string
	res := C.keygen_has_outgoing(k.state)
	k.trace("keygen_has_outgoing", res)
	for res > 0 {
		outgoingBytesSize := C.keygen_outgoing(k.state, (*C.char)(k.buffer), BufferSize)
		k.trace("keygen_outgoing_size", outgoingBytesSize)
		k.trace("keygen_outgoing", C.GoString((*C.char)(k.buffer)))
		outgoing = append(outgoing, C.GoString((*C.char)(k.buffer)))
		res = C.keygen_has_outgoing(k.state)
	}
	return outgoing
}

func (k *KeygenSimple) GetState() (*KeygenState, error) {

	res := C.keygen_get_state(k.state, (*C.char)(k.buffer), BufferSize)
	jsonStr := C.GoString((*C.char)(k.buffer))

	if res > 0 {
		k.trace("keygen_get_state", jsonStr)
		var state KeygenState
		err := json.Unmarshal([]byte(jsonStr), &state)
		if err != nil {
			return nil, err
		}
		return &state, nil
	}
	return nil, nil
}

func (k *KeygenSimple) handleIncoming(msg string) {
	k.trace("incoming", msg)
	cText := C.CString(msg)
	defer C.free(unsafe.Pointer(cText))
	C.keygen_incoming(k.state, cText)
}

func (k *KeygenSimple) finishIfPossible() *string {
	finished := C.keygen_is_finished(k.state)
	if finished != 1 {
		return nil
	}
	res := C.keygen_pick_output(k.state, (*C.char)(k.buffer), BufferSize)
	output := C.GoString((*C.char)(k.buffer))

	if res > 0 {
		k.trace("keygen_pick_output", output)
		return &output
	}
	return nil
}

func (k *KeygenSimple) trace(funcName string, result interface{}) {
	log.WithFields(log.Fields{
		"participant": k.i,
		"funcName":    funcName,
		"result":      result,
	}).Trace("statusCheck")
}
