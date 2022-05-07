package main

/*
#cgo LDFLAGS: -L./target/release -lmulti_party_bls_wrapper
#include <stdio.h>
#include <stdlib.h>
#include "./bls.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

const BufferSize = 2048

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
	fmt.Printf("%v keygen_wants_to_proceed: %v\n", k.i, res)
	if res == 1 {
		res = C.keygen_proceed(k.state)
		fmt.Printf("%v keygen_proceed: %v\n", k.i, res)
	}
}

func (k *Keygen) sendOutgoingIfThereIs() {
	res := C.keygen_has_outgoing(k.state)
	fmt.Printf("%v keygen_has_outgoing: %v\n", k.i, res)
	for res > 0 {
		outgoingBytesSize := C.keygen_outgoing(k.state, (*C.char)(k.buffer), BufferSize)

		fmt.Printf("%v outgoing bytes size: %v\n", k.i, outgoingBytesSize)
		fmt.Printf("%v outgoing is:\n", k.i)
		fmt.Printf("\033[0;32m")
		fmt.Printf("%s\n", C.GoString((*C.char)(k.buffer)))
		fmt.Printf("\033[0m")
		k.outgoing <- C.GoString((*C.char)(k.buffer))
		res = C.keygen_has_outgoing(k.state)
	}
}

//func (k *Keygen) waitForIncoming(msg string) {
//	fmt.Printf("incoming > ")
//	//fgets(buffer, BUFFER_SIZE, stdin);
//	//text, _ := stdin.ReadString('\n')
//	cText := C.CString(msg)
//	defer C.free(unsafe.Pointer(cText))
//	C.keygen_incoming(k.state, cText)
//}

func (k *Keygen) handleIncoming(msg string) {
	fmt.Printf("%v has incoming: %v\n", k.i, msg)
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
		fmt.Printf("%v Output is:\n%v\n", k.i, output)
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
			fmt.Printf("%v finished: %v\n", k.i, finished)
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

func (k *Keygen) Initialize() {
	k.proceedIfNeeded()
	//k.sendOutgoingIfThereIs()
	//k.finishIfPossible()
}

type ProtocolMessage struct {
	Sender   int                    `json:"sender"`
	Receiver int                    `json:"receiver"`
	Ignored  map[string]interface{} `json:"-"` // Rest of the fields should go here.
}

func main() {

	t := 1
	n := 3
	var (
		ins  []chan string
		outs []chan string
		machines    [] *Keygen
	)


	for i := 1; i < n + 1; i++ {
		in := make(chan string, n)
		out := make(chan string, n)
		keygen := NewKeygen(i, t, n, in, out)
		ins = append(ins, in)
		outs = append(outs, out)
		machines = append(machines, keygen)
	}

	defer func(machines []*Keygen){
		for _, machine := range machines {
			machine.Free()
		}
	}(machines)


	go func(o1 <-chan string, o2 <-chan string, o3 <-chan string, i1 chan<- string,i2 chan<- string,i3 chan<- string) {
		send := func (str string){
			msg := ProtocolMessage{}
			if err := json.Unmarshal([]byte(str), &msg); err != nil {
				fmt.Printf("error: %v\n", err)
			}else {
				switch msg.Receiver {
				case 0:
					if msg.Sender != 1 {
						i1 <- str
					}
					if msg.Sender != 2 {
						i2 <- str
					}
					if msg.Sender != 3 {
						i3 <- str
					}
				case 1:
					i1 <- str
				case 2:
					i2 <- str
				case 3:
					i3 <- str
				}
			}
		}
		for {
			select {
			case str, ok := <-o1:
				if ok {
					send(str)
				}
			case str, ok := <-o2:
				if ok {
					send(str)
				}
			case str, ok := <-o3:
				if ok {
					send(str)
				}
			case <-time.After(1 * time.Second):
			}
		}
	}(outs[0],outs[1],outs[2],ins[0],ins[1],ins[2])




	go machines[0].ProcessLoop()
	go machines[1].ProcessLoop()
	go machines[2].ProcessLoop()

	machines[0].Initialize()
	machines[1].Initialize()
	machines[2].Initialize()

	var allFinished bool
	for !allFinished {
		select {
		case <-time.After(1 * time.Second):
			allFinished = true
			for _, machine := range machines {
				allFinished = allFinished && machine.output != nil
			}
			fmt.Printf("allFinished: %v\n", allFinished)
			if allFinished {
				break
			}
		}
	}
}
