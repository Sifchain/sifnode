package legacy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DistributionRecord084 struct {
	ClaimStatus                 int64
	DistributionName            string         `json:"distribution_name"`
	RecipientAddress            sdk.AccAddress `json:"recipient_address"`
	Coins                       sdk.Coins      `json:"coins"`
	DistributionStartHeight     int64          `json:"distribution_start_height"`
	DistributionCompletedHeight int64          `json:"distribution_completed_height"`
}
