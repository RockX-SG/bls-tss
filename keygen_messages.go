package bls_tss

type KeygenRoundMsg struct {
	Sender   int         `json:"sender"`
	Receiver interface{} `json:"receiver"`
	Body     struct {
		Round1 *KeygenRound1 `json:"Round1,omitempty"`
		Round2 *KeygenRound2 `json:"Round2,omitempty"`
		Round3 *KeygenRound3 `json:"Round3,omitempty"`
		Round4 *KeygenRound4 `json:"Round4,omitempty"`
	} `json:"body"`
}

type KeygenRound1 struct {
	Com []int `json:"com"`
}

type KeygenRound2 struct {
	BlindFactor []int `json:"blind_factor"`
	YI          struct {
		Curve string `json:"curve"`
		Point []int  `json:"point"`
	} `json:"y_i"`
}

type KeygenRound3 struct {
	I           int `json:"i"`
	T           int `json:"t"`
	N           int `json:"n"`
	J           int `json:"j"`
	Commitments []struct {
		Curve string `json:"curve"`
		Point []int  `json:"point"`
	} `json:"commitments"`
	Share struct {
		Curve  string `json:"curve"`
		Scalar []int  `json:"scalar"`
	} `json:"share"`
}

type KeygenRound4 struct {
	Pk struct {
		Curve string `json:"curve"`
		Point []int  `json:"point"`
	} `json:"pk"`
	PkTRandCommitment struct {
		Curve string `json:"curve"`
		Point []int  `json:"point"`
	} `json:"pk_t_rand_commitment"`
	ChallengeResponse struct {
		Curve  string `json:"curve"`
		Scalar []int  `json:"scalar"`
	} `json:"challenge_response"`
}

type LocalKey struct {
	SharedKey struct {
		I  int `json:"i"`
		T  int `json:"t"`
		N  int `json:"n"`
		Vk struct {
			Curve string `json:"curve"`
			Point []int  `json:"point"`
		} `json:"vk"`
		SkI struct {
			Curve  string `json:"curve"`
			Scalar []int  `json:"scalar"`
		} `json:"sk_i"`
	} `json:"shared_key"`
	VkVec []struct {
		Curve string `json:"curve"`
		Point []int  `json:"point"`
	} `json:"vk_vec"`
}
