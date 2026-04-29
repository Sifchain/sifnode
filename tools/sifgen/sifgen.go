package sifgen

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MakeNowJust/heredoc"
	"github.com/serfersac/sifnode/gcp" // Import the GCP provider
	sifgentypes "github.com/serfersac/sifnode/tools/sifgen/common/types" // Import the common types
	"github.com/Sifchain/sifnode/tools/sifgen/key"
	"github.com/Sifchain/sifnode/tools/sifgen/network"
	"github.com/Sifchain/sifnode/tools/sifgen/node"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
)

type Sifgen struct {
	chainID *string
	providers map[string]sifgentypes.Provider
}

func NewSifgen(chainID *string) Sifgen {
	s := Sifgen{
		chainID: chainID,
		providers: make(map[string]sifgentypes.Provider),
	}
	// Register GCP provider
	ctx := context.Background() // Use a background context for initial provider setup
	gcpProvider, err := gcp.NewGCPProvider(ctx, "sifchain-devnet-307421", "asia-east1-a") // TODO: Make projectID and zone configurable
	if err != nil {
		log.Fatalf("Failed to create GCP provider: %v", err)
	}
	s.RegisterProvider("gcp", gcpProvider)
	return s
}

// RegisterProvider registers a new cloud provider.
func (s *Sifgen) RegisterProvider(name string, provider sifgentypes.Provider) {
	s.providers[name] = provider
}

func (s Sifgen) NewNetwork(keyringBackend string) *network.Network {
	return &network.Network{
		ChainID: *s.chainID,
		CLI:     utils.NewCLI(*s.chainID, keyringBackend),
	}
}

func (s Sifgen) NetworkCreate(ctx context.Context, count int, outputDir, startingIPAddress, providerName, outputFile string) {
	provider, ok := s.providers[providerName]
	if !ok {
		log.Fatalf("Provider '%s' not found", providerName)
	}

	net := network.NewNetwork(*s.chainID)
	// Assuming net.Build can take a provider to create nodes
	// This part needs further integration with the existing network.Build logic
	summary, err := net.Build(ctx, count, outputDir, startingIPAddress, provider)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = ioutil.WriteFile(outputFile, []byte(*summary), 0600); err != nil {
		log.Fatal(err)
		return
	}
}

func (s Sifgen) NetworkReset(networkDir string) {
	if err := network.Reset(*s.chainID, networkDir); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) NewNode(keyringBackend string) *node.Node {
	return &node.Node{
		ChainID: *s.chainID,
		CLI:     utils.NewCLI(*s.chainID, keyringBackend),
	}
}

func (s Sifgen) NodeReset(nodeHomeDir *string) {
	if err := node.Reset(*s.chainID, nodeHomeDir); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) KeyGenerateMnemonic(name, password string) {
	newKey := key.NewKey(name, password)
	newKey.GenerateMnemonic()
	fmt.Println(newKey.Mnemonic)
}

func (s Sifgen) KeyRecoverFromMnemonic(mnemonic string) {
	newKey := key.NewKey("", "")
	if err := newKey.RecoverFromMnemonic(mnemonic); err != nil {
		log.Fatal(err)
	}

	fmt.Println(heredoc.Doc(`
		Address: ` + newKey.Address + `
		Validator Address: ` + newKey.ValidatorAddress + `
		Consensus Address: ` + newKey.ConsensusAddress + `
	`))
}
