package v39

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "dispensation"

type DistributionStatus int64

const Pending DistributionStatus = 1
const Completed DistributionStatus = 2

func (ds DistributionStatus) String() string {
	switch ds {
	case Pending:
		return "Pending"
	case Completed:
		return "Completed"
	default:
		return "All"
	}
}

type GenesisState struct {
	DistributionRecords DistributionRecords `json:"distribution_records"`
	Distributions       Distributions       `json:"distributions"`
	Claims              UserClaims          `json:"claims"`
}

type DistributionRecord struct {
	DistributionStatus          DistributionStatus `json:"distribution_status"`
	DistributionName            string             `json:"distribution_name"`
	DistributionType            DistributionType   `json:"distribution_type"`
	RecipientAddress            sdk.AccAddress     `json:"recipient_address"`
	Coins                       sdk.Coins          `json:"coins"`
	DistributionStartHeight     int64              `json:"distribution_start_height"`
	DistributionCompletedHeight int64              `json:"distribution_completed_height"`
	AuthorizedRunner            sdk.AccAddress     `json:"authorized_runner"`
}

type DistributionRecords []DistributionRecord

type DistributionType int64

const DistributionTypeUnknown DistributionType = 0
const Airdrop DistributionType = 1
const LiquidityMining DistributionType = 2
const ValidatorSubsidy DistributionType = 3

func (dt DistributionType) String() string {
	switch dt {
	case Airdrop:
		return "Airdrop"
	case LiquidityMining:
		return "LiquidityMining"
	case ValidatorSubsidy:
		return "ValidatorSubsidy"
	default:
		return "Invalid"
	}
}

type Distribution struct {
	DistributionType DistributionType `json:"distribution_type"`
	DistributionName string           `json:"distribution_name"`
	Runner           sdk.AccAddress   `json:"runner"`
}

type Distributions []Distribution

type UserClaim struct {
	UserAddress   sdk.AccAddress   `json:"user_address"`
	UserClaimType DistributionType `json:"user_claim_type"`
	UserClaimTime time.Time        `json:"user_claim_time"`
}

type UserClaims []UserClaim
