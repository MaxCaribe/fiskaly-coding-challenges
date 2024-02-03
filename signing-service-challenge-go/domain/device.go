package domain

type SignatureDevice struct {
	UUID             string    `json:"uuid"`
	Label            string    `json:"label"`
	PrivateKey       []byte    `json:"private_key"`
	PublicKey        []byte    `json:"public_key"`
	Algorithm        Algorithm `json:"algorithm" `
	SignatureCounter int       `json:"signature_counter"`
	LastSignature    []byte    `json:"last_signature"`
}
