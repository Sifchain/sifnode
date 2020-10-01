package networks

import (
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/yelinaung/go-haikunator"
)

type Witness struct {
	name        string
	address     string
	peerAddress string
	keyPassword string
	genesisURL  string
	utils       Utils
}

func NewWitness(peerAddress, genesisURL, defaultNodeHome string) *Witness {
	keyPassword, _ := password.Generate(32, 5, 0, false, false)

	return &Witness{
		name:        haikunator.New(time.Now().UTC().UnixNano()).Haikunate(),
		peerAddress: peerAddress,
		keyPassword: keyPassword,
		genesisURL:  genesisURL,
		utils:       NewUtils(defaultNodeHome),
	}
}

func (w *Witness) Name() string {
	return w.name
}

func (w *Witness) Address(address *string) *string {
	if address == nil {
		return &w.address
	}

	w.address = *address
	return nil
}

func (w *Witness) PeerAddress() string {
	return w.peerAddress
}

func (w *Witness) KeyPassword() string {
	return w.keyPassword
}

func (w *Witness) GenesisURL() string {
	return w.genesisURL
}

func (w *Witness) CollectPeerAddress() {}
