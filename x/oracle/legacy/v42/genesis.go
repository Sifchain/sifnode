package v42

import (
	v039oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v39"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func Migrate(genesis v039oracle.GenesisState) *types.GenesisState {
	networkDescriptor := types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	whitelist := make(map[string]uint32)
	defaultPower := uint32(100)

	for _, addr := range genesis.AddressWhitelist {
		whitelist[addr.String()] = defaultPower
	}

	addressWhitelist := make(map[uint32]*types.ValidatorWhiteList)
	addressWhitelist[uint32(networkDescriptor)] = &types.ValidatorWhiteList{WhiteList: whitelist}

	return &types.GenesisState{
		AddressWhitelist: addressWhitelist,
		AdminAddress:     genesis.AdminAddress.String(),
		Prophecies:       []*types.Prophecy{},
	}
}
