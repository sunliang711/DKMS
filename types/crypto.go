package types

// CryptoConfig TODO
// 2019/08/17 12:03:28
type CryptoConfig struct {
	PID          string `json:"pid"`
	EccType      string `json:"ecc_type"`
	EncAlgorithm string `json:"enc_algorithm"`
	GenManner    string `json:"gen_manner"`
	T            int    `json:"t"`
	N            int    `json:"n"`
	API          string `json:"api"`
}
