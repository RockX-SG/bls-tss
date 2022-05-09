package main

import "C"
import (
	"encoding/json"
	"fmt"
	"time"
)


const BufferSize = 2048

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
