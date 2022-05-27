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
	"time"
	"unsafe"
)

type Sign struct {
	msgHash string
	i int
	n int
	localKey string
	state    unsafe.Pointer
	buffer   unsafe.Pointer
	incoming <-chan string
	outgoing chan<- string
	output *string
}

func NewSign(msgHash string, i int, n int, localKey string, incoming <-chan string, outgoing chan<- string) *Sign {
	cHash := C.CString(msgHash)
	defer C.free(unsafe.Pointer(cHash))
	cKey := C.CString(localKey)
	defer C.free(unsafe.Pointer(cKey))

	buffer := C.malloc(C.size_t(BufferSize))
	state := C.new_sign(cHash, C.int(i), C.int(n), cKey)
	// void* new_sign(const char* msg_hash, int i, int n, const char* local_key);
	return &Sign{msgHash, i, n, localKey, state, buffer, incoming, outgoing, nil}
}

func (k *Sign) Free() {
	C.free(k.buffer)
	C.free_sign(k.state)
}

func (k *Sign) Output() * string{
	return k.output
}

func (k *Sign) proceedIfNeeded() {
	res := C.sign_wants_to_proceed(k.state)
	k.trace("sign_wants_to_proceed", res)
	if res == 1 {
		res = C.sign_proceed(k.state)
		k.trace("sign_proceed", res)
	}
}

func (k *Sign) sendOutgoingIfThereIs() {
	res := C.sign_has_outgoing(k.state)
	k.trace("sign_has_outgoing", res)
	for res > 0 {
		outgoingBytesSize := C.sign_outgoing(k.state, (*C.char)(k.buffer), BufferSize)
		k.trace("sign_outgoing_size", outgoingBytesSize)
		k.trace("sign_outgoing", C.GoString((*C.char)(k.buffer)))
		k.outgoing <- C.GoString((*C.char)(k.buffer))
		res = C.sign_has_outgoing(k.state)
	}
}

//func (k *Sign) waitForIncoming(msg string) {
//	log.Debugf("incoming > ")
//	//fgets(buffer, BUFFER_SIZE, stdin);
//	//text, _ := stdin.ReadString('\n')
//	cText := C.CString(msg)
//	defer C.free(unsafe.Pointer(cText))
//	C.sign_incoming(k.state, cText)
//}

func (k *Sign) handleIncoming(msg string) {
	k.trace("incoming", msg)
	cText := C.CString(msg)
	defer C.free(unsafe.Pointer(cText))
	C.sign_incoming(k.state, cText)
}

func (k *Sign) finishIfPossible() {
	finished := C.sign_is_finished(k.state)
	if finished != 1 {
		return
	}
	res := C.sign_pick_output(k.state, (*C.char)(k.buffer), BufferSize)
	output := C.GoString((*C.char)(k.buffer))
	k.output = &output
	if res > 0 {
		k.trace("sign_pick_output", output)
	}
}

func (k *Sign) ProcessLoop() {
	var finished bool
	for !finished {
		select {
		case msg, ok := <-k.incoming:
			if ok {
				k.handleIncoming(msg)
				//k.sendOutgoingIfThereIs()
				//k.proceedIfNeeded()
				//k.sendOutgoingIfThereIs()
				//k.finishIfPossible()
			}
		case <-time.After(1 * time.Second):
			finished = k.output != nil
			k.trace("finished", finished)
			if finished {
				break
			} else {
				k.sendOutgoingIfThereIs()
				k.proceedIfNeeded()
				k.finishIfPossible()
			}
		}
	}
}

func (k *Sign) trace(funcName string, result interface{}){
	log.WithFields(log.Fields{
		"participant": k.i,
		"funcName": funcName,
		"result": result,
	}).Trace("statusCheck")
}

func (k *Sign) Initialize() {
	k.proceedIfNeeded()
	//k.sendOutgoingIfThereIs()
	//k.finishIfPossible()
}
