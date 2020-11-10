package types

type Nodes []struct {
	ChainID                   string `yaml:"chain_id"`
	NodeID                    string `yaml:"node_id"`
	Ipv4Address               string `yaml:"ipv4_address"`
	Moniker                   string `yaml:"moniker"`
	Password                  string `yaml:"password"`
	Address                   string `yaml:"address"`
	PubKey                    string `yaml:"pub_key"`
	ValidatorAddress          string `yaml:"validator_address"`
	ValidatorConsensusAddress string `yaml:"validator_consensus_address"`
	IsSeed                    bool   `yaml:"is_seed"`
}
