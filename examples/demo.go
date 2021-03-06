package main

import "C"
import (
	"encoding/json"
	"fmt"
	tss "github.com/RockX-SG/bls-tss"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})

	//Set log output to standard output (default output is stderr, standard error)
	//Log message output can be any io.writer type
	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

func main() {

	t := 1
	n := 3
	var (
		ins       []chan string
		outs      []chan string
		kMachines []*tss.Keygen
		sMachines []*tss.Sign
	)

	for i := 1; i < n+1; i++ {
		in := make(chan string, n)
		out := make(chan string, n)
		keygen := tss.NewKeygen(i, t, n, in, out)
		ins = append(ins, in)
		outs = append(outs, out)
		kMachines = append(kMachines, keygen)
	}

	defer func(machines []*tss.Keygen) {
		for _, machine := range machines {
			machine.Free()
		}
	}(kMachines)

	go func(o1 <-chan string, o2 <-chan string, o3 <-chan string, i1 chan<- string, i2 chan<- string, i3 chan<- string) {
		send := func(str string) {
			msg := tss.ProtocolMessage{}
			if err := json.Unmarshal([]byte(str), &msg); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
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
	}(outs[0], outs[1], outs[2], ins[0], ins[1], ins[2])

	log.Debug("Starting keygen")
	go kMachines[0].ProcessLoop()
	go kMachines[1].ProcessLoop()
	go kMachines[2].ProcessLoop()

	kMachines[0].Initialize()
	kMachines[1].Initialize()
	kMachines[2].Initialize()
	log.Debug("KeygenSimple started")

	var allFinished bool
	for !allFinished {
		select {
		case <-time.After(1 * time.Second):
			allFinished = true
			for _, machine := range kMachines {
				allFinished = allFinished && machine.Output() != nil
			}
			log.Tracef("keygen allFinished: %v\n", allFinished)
			if allFinished {
				break
			}
		}
	}
	log.Debug("KeygenSimple completed")
	for _, machine := range kMachines {
		log.WithFields(log.Fields{
			"result": *machine.Output(),
		}).Trace("keygen result")
	}

	msgHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 32 * "a"

	n = 2
	for i := 1; i < n+1; i++ {
		sign := tss.NewSign(msgHash, i, n, *kMachines[i-1].Output(), ins[i-1], outs[i-1])
		sMachines = append(sMachines, sign)
	}

	defer func(machines []*tss.Sign) {
		for _, machine := range machines {
			machine.Free()
		}
	}(sMachines)

	log.Debug("Start signing")
	go sMachines[0].ProcessLoop()
	go sMachines[1].ProcessLoop()

	sMachines[0].Initialize()
	sMachines[1].Initialize()
	//sMachines[2].Initialize()
	log.Debug("Signing started")

	allFinished = false
	for !allFinished {
		select {
		case <-time.After(1 * time.Second):
			allFinished = true
			for _, machine := range sMachines {
				allFinished = allFinished && machine.Output() != nil
			}
			log.Tracef("sign allFinished: %v\n", allFinished)
			if allFinished {
				break
			}
		}
	}
	log.Debug("Signing completed")
	log.Infof("msgHash is: %v\n", msgHash)
	log.Infof("signature is: %v\n", *sMachines[1].Output())

}
