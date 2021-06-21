package v42

import (
	v039oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v39"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func Migrate(genesis v039oracle.GenesisState) *types.GenesisState {
	var addressWhiteList = make([]string, len(genesis.AddressWhitelist))
	for i, addr := range genesis.AddressWhitelist {
		addressWhiteList[i] = addr.String()
	}

	return &types.GenesisState{
		AddressWhitelist: addressWhiteList,
		AdminAddress:     genesis.AdminAddress.String(),
		// TODO: Add prophecies once defined in 39&42 genesis state
	}
}
