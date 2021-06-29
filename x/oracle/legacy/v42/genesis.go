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

	prophecies := make([]*types.DBProphecy, len(genesis.Prophecies))
	for index, legacy := range genesis.Prophecies {

		statusText := types.StatusText_STATUS_TEXT_UNSPECIFIED
		if legacy.Status.Text == v039oracle.PendingStatusText {
			statusText = types.StatusText_STATUS_TEXT_PENDING
		} else if legacy.Status.Text == v039oracle.FailedStatusText {
			statusText = types.StatusText_STATUS_TEXT_FAILED
		} else if legacy.Status.Text == v039oracle.SuccessStatusText {
			statusText = types.StatusText_STATUS_TEXT_SUCCESS
		}

		prophecies[index] = &types.DBProphecy{
			Id: legacy.ID,
			Status: types.Status{
				Text:       statusText,
				FinalClaim: legacy.Status.FinalClaim,
			},
			ClaimValidators: legacy.ClaimValidators,
			ValidatorClaims: legacy.ValidatorClaims,
		}
	}
	addressWhitelist := make(map[uint32]*types.ValidatorWhiteList)
	addressWhitelist[uint32(networkDescriptor)] = &types.ValidatorWhiteList{WhiteList: whitelist}

	return &types.GenesisState{
		AddressWhitelist: addressWhitelist,
		AdminAddress:     genesis.AdminAddress.String(),
		Prophecies:       prophecies,
	}
}
