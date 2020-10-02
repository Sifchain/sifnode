package networks

import (
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/yelinaung/go-haikunator"
)

type Witness struct {
	moniker     string
	address     string
	peerAddress string
	keyPassword string
	genesisURL  string
	utils       NetworkUtils
}

func NewWitness(peerAddress, genesisURL string, utils NetworkUtils) *Witness {
	keyPassword, _ := password.Generate(32, 5, 0, false, false)

	return &Witness{
		moniker:     haikunator.New(time.Now().UTC().UnixNano()).Haikunate(),
		peerAddress: peerAddress,
		keyPassword: keyPassword,
		genesisURL:  genesisURL,
		utils:       utils,
	}
}

func (w *Witness) Moniker() string {
	return w.moniker
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

func (w *Witness) CollectPeerAddress() error { return nil }
