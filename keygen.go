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

type Keygen struct {
	i int
	t int
	n int
	state    unsafe.Pointer
	buffer   unsafe.Pointer
	incoming <-chan string
	outgoing chan<- string
	output *string
}

func NewKeygen(i int, t int, n int, incoming <-chan string, outgoing chan<- string) *Keygen {
	buffer := C.malloc(C.size_t(BufferSize))
	state := C.new_keygen(C.int(i), C.int(t), C.int(n))
	return &Keygen{i, t, n, state, buffer, incoming, outgoing, nil}
}

func (k *Keygen) Free() {
	C.free(k.buffer)
	C.free_keygen(k.state)
}

func (k *Keygen) Output() * string{
	return k.output
}

func (k *Keygen) proceedIfNeeded() {
	res := C.keygen_wants_to_proceed(k.state)
	k.trace("keygen_wants_to_proceed", res)
	if res == 1 {
		res = C.keygen_proceed(k.state)
		k.trace("keygen_proceed", res)
	}
}

func (k *Keygen) sendOutgoingIfThereIs() {
	res := C.keygen_has_outgoing(k.state)
	k.trace("keygen_has_outgoing", res)
	for res > 0 {
		outgoingBytesSize := C.keygen_outgoing(k.state, (*C.char)(k.buffer), BufferSize)
		k.trace("keygen_outgoing_size", outgoingBytesSize)
		k.trace("keygen_outgoing", C.GoString((*C.char)(k.buffer)))
		k.outgoing <- C.GoString((*C.char)(k.buffer))
		res = C.keygen_has_outgoing(k.state)
	}
}

//func (k *Keygen) waitForIncoming(msg string) {
//	log.Debugf("incoming > ")
//	//fgets(buffer, BUFFER_SIZE, stdin);
//	//text, _ := stdin.ReadString('\n')
//	cText := C.CString(msg)
//	defer C.free(unsafe.Pointer(cText))
//	C.keygen_incoming(k.state, cText)
//}

func (k *Keygen) handleIncoming(msg string) {
	k.trace("incoming", msg)
	cText := C.CString(msg)
	defer C.free(unsafe.Pointer(cText))
	C.keygen_incoming(k.state, cText)
}

func (k *Keygen) finishIfPossible() {
	finished := C.keygen_is_finished(k.state)
	if finished != 1 {
		return
	}
	res := C.keygen_pick_output(k.state, (*C.char)(k.buffer), BufferSize)
	output := C.GoString((*C.char)(k.buffer))
	k.output = &output
	if res > 0 {
		k.trace("keygen_pick_output", output)
	}
}

func (k *Keygen) ProcessLoop() {
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

func (k *Keygen) trace(funcName string, result interface{}){
	log.WithFields(log.Fields{
		"participant": k.i,
		"funcName": funcName,
		"result": result,
	}).Trace("statusCheck")
}

func (k *Keygen) Initialize() {
	k.proceedIfNeeded()
	//k.sendOutgoingIfThereIs()
	//k.finishIfPossible()
}
