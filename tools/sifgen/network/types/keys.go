package types

type NetworkKeys []Keys
type Keys []Key
type Key struct {
	NodeID   string `json:"node_id"`
	Name     string `json:"name" yaml:"name"`
	KeyType  string `json:"type" yaml:"type"`
	Address  string `json:"address" yaml:"address"`
	PubKey   string `json:"pubkey" yaml:"pubkey"`
	Moniker  string `json:"moniker"`
	Password string `json:"password"`
}
