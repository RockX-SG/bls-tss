package bls_tss

type SignRoundMsg struct {
	Sender   int         `json:"sender"`
	Receiver interface{} `json:"receiver"`
	Body     struct {
		Round1 *SignRound1 `json:"Round1,omitempty"`
	} `json:"body"`
}

type SignRound1 = PartialSignature

type PartialSignature struct {
	I      int `json:"i"`
	SigmaI struct {
		Curve string `json:"curve"`
		Point []int  `json:"point"`
	} `json:"sigma_i"`
	DdhProof struct {
		A1 struct {
			Curve string `json:"curve"`
			Point []int  `json:"point"`
		} `json:"a1"`
		A2 struct {
			Curve string `json:"curve"`
			Point []int  `json:"point"`
		} `json:"a2"`
		Z []int `json:"z"`
	} `json:"ddh_proof"`
}
