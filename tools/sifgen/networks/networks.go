package networks

var (
	Coins = []string{
		"1000000000stake",
		"100000rowan",
	}
)

type Network interface {
	Reset()
	Setup()
	Genesis()
}

type NetworkNode interface {
	Name() string
	Address(*string) *string
	PeerAddress() string
	KeyPassword() string
	GenesisURL() string
	CollectPeerAddress()
}
