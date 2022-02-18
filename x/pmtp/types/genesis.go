package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: &Params{},
	}
}
