package bls_tss

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeygenRoundMsg(t *testing.T) {
	examples := []string{
		"{\"sender\":1,\"receiver\":null,\"body\":{\"Round1\":{\"com\":[185,120,129,10,16,104,222,69,34,233,113,207,109,235,201,214,62,31,179,215,155,224,22,32,22,49,23,94,239,187,50,141]}}}",
		"{\"sender\":2,\"receiver\":null,\"body\":{\"Round2\":{\"blind_factor\":[51,120,251,107,20,146,180,151,203,35,235,96,111,50,195,210,45,114,22,104,0,85,25,228,47,207,48,205,164,83,115,58],\"y_i\":{\"curve\":\"bls12_381_1\",\"point\":[164,128,21,238,237,122,161,90,158,231,169,112,135,193,90,83,211,144,44,64,247,65,137,181,230,129,53,59,241,167,114,159,165,64,26,242,229,199,198,174,119,101,15,103,231,229,118,164]}}}}",
		"{\"sender\":1,\"receiver\":2,\"body\":{\"Round3\":{\"i\":1,\"t\":1,\"n\":3,\"j\":2,\"commitments\":[{\"curve\":\"bls12_381_1\",\"point\":[164,199,103,95,33,143,94,62,79,117,51,160,28,112,202,157,97,178,237,202,144,134,146,234,64,110,78,136,115,151,83,21,150,22,186,33,251,56,95,181,12,202,209,51,50,177,66,37]},{\"curve\":\"bls12_381_1\",\"point\":[177,36,18,86,162,6,186,210,239,97,124,57,196,215,2,124,28,43,248,231,126,134,20,109,249,21,50,6,255,137,170,18,153,107,30,39,113,13,52,16,165,160,109,245,239,244,100,138]}],\"share\":{\"curve\":\"bls12_381_1\",\"scalar\":[22,174,122,230,245,66,201,4,116,207,36,43,122,185,150,255,59,113,3,63,243,188,43,24,189,39,231,78,93,88,10,55]}}}}",
		"{\"sender\":3,\"receiver\":null,\"body\":{\"Round4\":{\"pk\":{\"curve\":\"bls12_381_1\",\"point\":[169,108,148,58,128,112,173,237,95,247,127,68,74,153,99,110,252,80,104,254,204,126,61,125,240,171,249,249,67,45,71,174,183,64,68,133,0,141,227,7,178,181,81,195,216,43,42,151]},\"pk_t_rand_commitment\":{\"curve\":\"bls12_381_1\",\"point\":[142,194,149,130,85,62,216,156,254,10,71,118,237,141,180,93,246,137,63,153,129,149,177,199,122,14,2,10,47,115,87,194,92,171,60,197,135,52,114,216,227,227,148,36,53,51,77,81]},\"challenge_response\":{\"curve\":\"bls12_381_1\",\"scalar\":[99,104,191,77,6,127,229,165,109,70,146,59,70,109,10,59,87,1,124,50,50,130,87,9,56,174,164,66,216,172,77,103]}}}}",
	}
	var err error
	for _, example := range examples {
		msg := KeygenRoundMsg{}
		err = json.Unmarshal([]byte(example), &msg)
		assert.Nil(t, err)
		jsonBytes, err := json.Marshal(msg)
		assert.Nil(t, err)
		assert.Equal(t, example, string(jsonBytes))
	}
}

func TestKeygenState(t *testing.T) {
	examples := []string{
		`{"Round0":{"party_i":1,"t":2,"n":4}}`,
		`{"Round1":{"keys":{"u_i":{"curve":"bls12_381_1","scalar":[209,203,234,225,7,138,31,40,218,72,193,42,57,36,148,54,193,172,49,117,144,32,96,66,222,245,84,54,221,119,142,86]},"y_i":{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]},"i":1},"comm":{"com":[221,4,42,72,177,12,142,114,244,121,84,232,16,157,52,0,219,94,70,175,62,106,147,121,223,3,4,116,29,57,60,130]},"decom":{"blind_factor":[203,94,207,107,115,34,21,248,158,96,177,252,143,234,207,215,163,192,137,72,43,198,76,113,232,39,44,127,203,129,242,15],"y_i":{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]}},"party_i":1,"t":2,"n":4}}`,
		`{"Round2":{"keys":{"u_i":{"curve":"bls12_381_1","scalar":[209,203,234,225,7,138,31,40,218,72,193,42,57,36,148,54,193,172,49,117,144,32,96,66,222,245,84,54,221,119,142,86]},"y_i":{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]},"i":1},"received_comm":[{"com":[221,4,42,72,177,12,142,114,244,121,84,232,16,157,52,0,219,94,70,175,62,106,147,121,223,3,4,116,29,57,60,130]},{"com":[237,71,170,85,167,37,123,129,64,154,222,34,218,242,24,64,187,157,170,12,251,229,182,157,37,4,195,165,19,100,214,160]},{"com":[183,74,227,39,45,132,144,191,177,101,194,207,55,218,109,226,171,224,75,81,223,230,180,149,153,84,14,82,108,86,165,195]},{"com":[6,27,38,103,155,202,137,84,110,102,108,253,254,90,70,70,94,126,111,23,104,225,130,55,90,29,51,175,245,175,114,110]}],"decom":{"blind_factor":[203,94,207,107,115,34,21,248,158,96,177,252,143,234,207,215,163,192,137,72,43,198,76,113,232,39,44,127,203,129,242,15],"y_i":{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]}},"party_i":1,"t":2,"n":4}}`,
		`{"Round3":{"keys":{"u_i":{"curve":"bls12_381_1","scalar":[209,203,234,225,7,138,31,40,218,72,193,42,57,36,148,54,193,172,49,117,144,32,96,66,222,245,84,54,221,119,142,86]},"y_i":{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]},"i":1},"y_vec":[{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]},{"curve":"bls12_381_1","point":[142,79,45,107,155,174,63,211,103,164,29,26,251,10,116,177,53,62,48,168,157,192,246,213,52,135,201,224,171,60,248,25,177,2,58,35,114,203,9,1,0,0,29,141,52,133,191,66]},{"curve":"bls12_381_1","point":[143,61,210,243,78,135,19,97,67,73,253,182,11,145,53,7,11,16,6,199,9,73,23,148,28,248,224,162,129,115,192,225,190,40,35,124,66,238,209,252,222,247,113,24,153,183,215,194]},{"curve":"bls12_381_1","point":[170,241,191,109,53,207,175,252,9,183,184,163,229,64,27,64,54,64,79,97,146,63,0,75,4,84,171,101,188,249,232,125,55,6,131,73,33,186,214,76,127,145,10,103,207,163,203,198]}],"own_share":{"i":1,"t":2,"n":4,"j":1,"commitments":[{"curve":"bls12_381_1","point":[177,12,55,133,210,20,208,54,22,234,205,90,83,200,167,205,135,70,15,92,39,174,216,23,24,156,36,201,245,0,80,50,97,190,41,98,88,66,28,199,191,147,58,67,26,73,76,210]},{"curve":"bls12_381_1","point":[153,67,203,128,95,158,156,55,164,135,29,213,158,107,119,62,195,66,224,137,138,252,25,70,175,245,233,24,69,113,239,215,252,147,210,26,84,190,253,201,227,159,242,37,167,103,159,84]},{"curve":"bls12_381_1","point":[128,1,63,146,116,245,101,200,163,99,56,50,202,7,196,22,191,126,29,220,226,253,64,124,53,249,206,197,198,244,124,127,199,68,34,224,212,54,249,241,103,50,94,45,42,147,92,25]}],"share":{"curve":"bls12_381_1","scalar":[35,201,45,189,127,120,204,112,124,57,86,68,97,154,169,96,133,130,195,40,1,145,137,27,161,248,243,104,242,88,186,59]}},"party_i":1,"t":2,"n":4}}`,
		`{"Round4":{"shared_keys":{"i":1,"t":2,"n":4,"vk":{"curve":"bls12_381_1","point":[172,8,232,70,161,53,195,211,250,90,61,160,136,186,138,18,241,16,214,64,181,208,148,239,156,97,150,156,5,58,231,97,40,61,247,2,22,91,244,229,78,16,21,35,204,18,66,169]},"sk_i":{"curve":"bls12_381_1","scalar":[81,88,207,4,11,55,125,69,246,248,200,159,61,224,188,125,250,24,45,123,23,185,80,222,42,67,129,192,22,100,13,43]}},"own_dlog_proof":{"pk":{"curve":"bls12_381_1","point":[141,195,161,3,79,63,63,52,12,81,243,5,57,69,203,149,220,249,65,31,231,185,43,209,106,36,168,103,205,148,197,13,26,182,74,185,106,156,246,119,134,189,60,194,79,136,26,39]},"pk_t_rand_commitment":{"curve":"bls12_381_1","point":[132,41,176,135,57,150,196,14,59,210,251,185,207,139,254,67,187,217,227,136,235,179,248,231,198,154,106,31,210,169,11,235,147,245,160,173,168,253,53,242,138,163,238,186,169,174,26,134]},"challenge_response":{"curve":"bls12_381_1","scalar":[229,234,11,72,254,188,23,64,156,75,242,0,49,220,97,139,197,176,220,109,86,156,202,193,130,145,102,9,197,141,169,31]}},"party_i":1,"t":2,"n":4}}`,
	}
	var err error
	for _, example := range examples {
		msg := KeygenState{}
		err = json.Unmarshal([]byte(example), &msg)
		assert.Nil(t, err)
		jsonBytes, err := json.Marshal(msg)
		assert.Nil(t, err)
		assert.Equal(t, example, string(jsonBytes))
	}
}

func TestLocalKey(t *testing.T) {
	example := "{\"shared_key\":{\"i\":3,\"t\":1,\"n\":3,\"vk\":{\"curve\":\"bls12_381_1\",\"point\":[147,245,81,1,38,241,14,113,141,44,43,234,66,157,115,75,154,57,216,99,179,3,65,61,182,64,101,41,223,180,238,193,212,95,229,126,219,41,251,244,164,214,179,43,238,150,91,226]},\"sk_i\":{\"curve\":\"bls12_381_1\",\"scalar\":[204,7,171,107,137,1,48,214,94,148,179,255,218,56,213,211,22,193,139,174,74,240,18,37,143,164,91,217,218,189,247,44]}},\"vk_vec\":[{\"curve\":\"bls12_381_1\",\"point\":[128,77,56,128,159,197,88,101,57,220,103,42,145,96,225,134,39,101,49,131,50,204,59,200,180,225,238,20,220,24,199,194,192,84,163,216,46,170,53,125,234,26,155,130,77,34,130,144]},{\"curve\":\"bls12_381_1\",\"point\":[173,200,135,100,170,254,19,58,227,28,231,250,130,152,105,141,120,66,218,46,29,157,91,149,42,244,77,122,100,24,163,158,101,0,173,191,90,132,166,139,100,74,143,77,91,88,114,81]},{\"curve\":\"bls12_381_1\",\"point\":[169,108,148,58,128,112,173,237,95,247,127,68,74,153,99,110,252,80,104,254,204,126,61,125,240,171,249,249,67,45,71,174,183,64,68,133,0,141,227,7,178,181,81,195,216,43,42,151]}]}"
	msg := LocalKey{}
	err := json.Unmarshal([]byte(example), &msg)
	assert.Nil(t, err)
	jsonBytes, err := json.Marshal(msg)
	assert.Nil(t, err)
	assert.Equal(t, example, string(jsonBytes))
}