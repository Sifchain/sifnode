package v42

import (
	v039oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v39"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func Migrate(genesis v039oracle.GenesisState) *types.GenesisState {
	return &types.GenesisState{
		// for new peggy2, each validator has its voting power, can be got from peggy 1.0
		AddressWhitelist: map[uint32]*types.ValidatorWhiteList{},
		AdminAddress:     genesis.AdminAddress.String(),
		// the algorithm to compute the prophecy id changed, not make sense to copy prophecy from peggy 1.0
		Prophecies: []*types.Prophecy{},
	}
}
