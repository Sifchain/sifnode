package types

type Keys []Key
type Key struct {
	Name      string   `json:"name" yaml:"name"`
	KeyType   string   `json:"type" yaml:"type"`
	Address   string   `json:"address" yaml:"address"`
	PubKey    PubKey   `json:"pubkey" yaml:"pubkey"`
	Mnemonic  string   `json:"mnemonic" yaml:"mnemonic"`
	Threshold int      `json:"threshold" yaml:"threshold"`
	PubKeys   []string `json:"pubkeys" yaml:"pubkeys"`
}

type PubKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}
