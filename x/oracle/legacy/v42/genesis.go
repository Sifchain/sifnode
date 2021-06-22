package v42

import (
	v039oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v39"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func Migrate(genesis v039oracle.GenesisState) *types.GenesisState {
	var addressWhiteList []string = make([]string, len(genesis.AddressWhitelist))
	for i, addr := range genesis.AddressWhitelist {
		addressWhiteList[i] = addr.String()
	}

	var prophecies []*types.DBProphecy
	for _, legacy := range genesis.Prophecies {

		statusText := types.StatusText_STATUS_TEXT_UNSPECIFIED
		if legacy.Status.Text == v039oracle.PendingStatusText {
			statusText = types.StatusText_STATUS_TEXT_PENDING
		} else if legacy.Status.Text == v039oracle.FailedStatusText {
			statusText = types.StatusText_STATUS_TEXT_FAILED
		} else if legacy.Status.Text == v039oracle.SuccessStatusText {
			statusText = types.StatusText_STATUS_TEXT_SUCCESS
		}

		prophecies = append(prophecies, &types.DBProphecy{
			Id: legacy.ID,
			Status: types.Status{
				Text:       statusText,
				FinalClaim: legacy.Status.FinalClaim,
			},
			ClaimValidators: legacy.ClaimValidators,
			ValidatorClaims: legacy.ValidatorClaims,
		})
	}

	return &types.GenesisState{
		AddressWhitelist: addressWhiteList,
		AdminAddress:     genesis.AdminAddress.String(),
		Prophecies:       prophecies,
	}
}
