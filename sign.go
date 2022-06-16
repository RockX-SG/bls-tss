package bls_tss

import "C"
import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Sign struct {
	inner           *SignSimple
	incoming        <-chan string
	outgoing        chan<- string
	pendingOutgoing []string
}

func NewSign(msgHash string, i int, n int, localKey string, incoming <-chan string, outgoing chan<- string) *Sign {
	inner := NewSignSimple(msgHash, i, n, localKey)
	// void* new_sign(const char* msg_hash, int i, int n, const char* local_key);
	return &Sign{inner, incoming, outgoing, nil}
}

func (k *Sign) Free() {
	k.inner.Free()
}

func (k *Sign) Initialize() {
	outgoing := k.inner.Init()
	k.pendingOutgoing = append(k.pendingOutgoing, outgoing...)
}

func (k *Sign) Output() *string {
	return k.inner.Output()
}

func (k *Sign) ProcessLoop() {
	var finished bool
	for !finished {
		select {
		case msg, ok := <-k.incoming:
			if ok {
				_, outgoing, _ := k.inner.Handle(msg)
				k.pendingOutgoing = append(k.pendingOutgoing, outgoing...)
				k.trace("outgoing", len(outgoing))
			}
		case <-time.After(1 * time.Second):
			finished := k.inner.Output() != nil
			k.trace("finished", finished)
			k.sendOutgoingIfThereIs()
			if finished {
				break
			}
		}
	}
}

func (k *Sign) sendOutgoingIfThereIs() {
	for len(k.pendingOutgoing) > 0 {
		item := k.pendingOutgoing[0]
		k.pendingOutgoing = k.pendingOutgoing[1:]
		k.outgoing <- item
	}
}

func (k *Sign) trace(funcName string, result interface{}) {
	log.WithFields(log.Fields{
		"participant": k.inner.i,
		"funcName":    funcName,
		"result":      result,
	}).Trace("statusCheck")
}
