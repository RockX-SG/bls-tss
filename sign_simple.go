package bls_tss

/*
#cgo CFLAGS:-I${SRCDIR}/bls-tss/include
#include <stdio.h>
#include <stdlib.h>
#include <tss.h>
*/
import "C"
import (
	log "github.com/sirupsen/logrus"
	"unsafe"
)

type SignSimple struct {
	msgHash  string
	i        int
	n        int
	localKey string
	state    unsafe.Pointer
	buffer   unsafe.Pointer
	output *string
}

func NewSignSimple(msgHash string, i int, n int, localKey string) *SignSimple {
	cHash := C.CString(msgHash)
	defer C.free(unsafe.Pointer(cHash))
	cKey := C.CString(localKey)
	defer C.free(unsafe.Pointer(cKey))

	buffer := C.malloc(C.size_t(BufferSize))
	state := C.new_sign(cHash, C.int(i), C.int(n), cKey)
	return &SignSimple{msgHash, i, n, localKey, state, buffer, nil}
}

func (k *SignSimple) Free() {
	C.free(k.buffer)
	C.free_sign(k.state)
}

func (k *SignSimple) Init() []string {
	k.proceedIfNeeded()
	return k.getOutgoing()
}

func (k *SignSimple) Handle(msg string) (bool, []string, error) {
	k.handleIncoming(msg)
	k.proceedIfNeeded()
	outgoing := k.getOutgoing()
	output := k.finishIfPossible()
	if len(outgoing) > 0 {
		return false, outgoing, nil
	}
	if output != nil {
		k.output = output
		return true, []string{*output}, nil
	}
	return false, nil, nil
}

func (k *SignSimple) Output() *string {
	return k.output
}

func (k *SignSimple) proceedIfNeeded() {
	res := C.sign_wants_to_proceed(k.state)
	k.trace("sign_wants_to_proceed", res)
	if res == 1 {
		res = C.sign_proceed(k.state)
		k.trace("sign_proceed", res)
	}
}

func (k *SignSimple) getOutgoing() []string {
	var outgoing []string
	res := C.sign_has_outgoing(k.state)
	k.trace("sign_has_outgoing", res)
	for res > 0 {
		outgoingBytesSize := C.sign_outgoing(k.state, (*C.char)(k.buffer), BufferSize)
		k.trace("sign_outgoing_size", outgoingBytesSize)
		k.trace("sign_outgoing", C.GoString((*C.char)(k.buffer)))
		outgoing = append(outgoing, C.GoString((*C.char)(k.buffer)))
		res = C.sign_has_outgoing(k.state)
	}
	return outgoing
}

func (k *SignSimple) handleIncoming(msg string) {
	k.trace("incoming", msg)
	cText := C.CString(msg)
	defer C.free(unsafe.Pointer(cText))
	C.sign_incoming(k.state, cText)
}

func (k *SignSimple) finishIfPossible() *string {
	finished := C.sign_is_finished(k.state)
	if finished != 1 {
		return nil
	}
	res := C.sign_pick_output(k.state, (*C.char)(k.buffer), BufferSize)
	if res > 0 {
		output := C.GoString((*C.char)(k.buffer))
		k.trace("sign_pick_output", output)
		k.output = &output
		return &output
	}
	return nil
}

func (k *SignSimple) trace(funcName string, result interface{}) {
	log.WithFields(log.Fields{
		"participant": k.i,
		"funcName":    funcName,
		"result":      result,
	}).Trace("statusCheck")
}
