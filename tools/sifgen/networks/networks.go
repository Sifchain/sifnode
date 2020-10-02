package networks

import (
	"github.com/Sifchain/sifnode/tools/sifgen/networks/types"
)

var (
	Coins = []string{
		"1000000000stake",
		"100000rowan",
	}
)

// Network (localnet/testnet/mainnet).
type Network interface {
	Setup() error
	Genesis() error
}

// Network nodes (currently just supports validators and witnesses).
type NetworkNode interface {
	Name() string
	Address(*string) *string
	PeerAddress() string
	KeyPassword() string
	GenesisURL() string
	CollectPeerAddress() error
}

// Network utils performs all the underlying tasks.
type NetworkUtils interface {
	Reset([]string) error
	InitChain(string, string) (*string, error)
	SetKeyRingStorage() (*string, error)
	SetConfigChainID(string) (*string, error)
	SetConfigIndent(bool) (*string, error)
	SetConfigTrustNode(bool) (*string, error)
	AddKey(string, string) (*string, error)
	AddGenesisAccount(string, []string) (*string, error)
	GenerateGenesisTxn(string, string) (*string, error)
	CollectGenesisTxns() (*string, error)
	ExportGenesis() (*string, error)
	GenesisFilePath() string
	ConfigFilePath() string
	ScrapePeerGenesis(string) (types.Genesis, error)
	SaveGenesis(types.Genesis) error
	ReplacePeerConfig([]string) error
}
