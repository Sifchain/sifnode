package v42

import (
	v039oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v39"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func Migrate(genesis v039oracle.GenesisState) *types.GenesisState {
	var addressWhiteList []string
	for _, addr := range genesis.AddressWhitelist {
		addressWhiteList = append(addressWhiteList, addr.String())
	}

	state := &types.GenesisState{
		AddressWhitelist: addressWhiteList,
		AdminAddress:     genesis.AdminAddress.String(),
		// TODO: Add prophecies once defined in 39&42 genesis state
	}
	return state
}
