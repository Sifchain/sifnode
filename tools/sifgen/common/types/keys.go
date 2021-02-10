package types

type Keys []Key
type Key struct {
	Name    string `json:"name" yaml:"name"`
	KeyType string `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
	PubKey  string `json:"pubkey" yaml:"pubkey"`
}
