package bls_tss

const BufferSize = 2048

type ProtocolMessage struct {
	Sender   int                    `json:"sender"`
	Receiver int                    `json:"receiver"`
	Ignored  map[string]interface{} `json:"-"` // Rest of the fields should go here.
}
