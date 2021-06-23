package types

// GenesisState - all faucet state that must be provided at genesis
type GenesisState struct {
	// TODO: The amount of fund that a user withdrew
}

// TODO : Implement the genesis functions
// NewGenesisState creates a new GenesisState object
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis validates the faucet genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
