package legacy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/dispensation/types"
)

type DistributionRecord084 struct {
	ClaimStatus                 int64          `json:"ClaimStatus"`
	DistributionName            string         `json:"distribution_name"`
	RecipientAddress            sdk.AccAddress `json:"recipient_address"`
	Coins                       sdk.Coins      `json:"coins"`
	DistributionStartHeight     int64          `json:"distribution_start_height"`
	DistributionCompletedHeight int64          `json:"distribution_completed_height"`
}

type Distribution084 struct {
	DistributionType types.DistributionType `json:"distribution_type"`
	DistributionName string                 `json:"distribution_name"`
}
