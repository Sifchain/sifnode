package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

<<<<<<< HEAD
func NewDistributionRecord(status DistributionStatus, distributionType DistributionType, distributionName string, recipientAddress string, coins sdk.Coins, start int64, end int64) DistributionRecord {
	return DistributionRecord{DistributionStatus: status, DistributionType: distributionType, DistributionName: distributionName, RecipientAddress: recipientAddress, Coins: coins, DistributionStartHeight: start, DistributionCompletedHeight: end}
=======
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

func IsValidStatus(s string) (DistributionStatus, bool) {
	switch s {
	case "Pending":
		return Pending, true
	case "Completed":
		return Completed, true
	default:
		return -1, false
	}
}

// DistributionRecord is created for every recipient for a distribution
// TODO : Add DistributionStatus to the prefixed key for records
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

func (records DistributionRecords) String() string {
	if len(records) == 0 {
		return ""
	}
	var rc string
	for _, record := range records {
		rc = rc + record.RecipientAddress.String() + ","
	}
	rc = rc[:len(rc)-1]
	return rc
}

func NewDistributionRecord(distributionName string, distributionType DistributionType, recipientAddress sdk.AccAddress, coins sdk.Coins, start int64, end int64, authorizedRunner sdk.AccAddress) DistributionRecord {
	return DistributionRecord{
		DistributionName:            distributionName,
		DistributionType:            distributionType,
		RecipientAddress:            recipientAddress,
		Coins:                       coins,
		DistributionStartHeight:     start,
		DistributionCompletedHeight: end,
		AuthorizedRunner:            authorizedRunner}
}

func (dr DistributionRecord) DoesClaimExist() bool {
	if dr.DistributionType == LiquidityMining || dr.DistributionType == ValidatorSubsidy {
		return true
	}
	return false
>>>>>>> develop
}

func (dr DistributionRecord) Validate() bool {
	if len(dr.RecipientAddress) == 0 {
		return false
	}
	if !dr.Coins.IsValid() {
		return false
	}
	if !dr.Coins.IsAllPositive() {
		return false
	}
	_, err := sdk.AccAddressFromBech32(dr.AuthorizedRunner.String())
	if err != nil {
		return false
	}
	_, err = sdk.AccAddressFromBech32(dr.RecipientAddress.String())
	if err != nil {
		return false
	}
	_, ok := IsValidDistributionType(dr.DistributionType.String())
	if !ok {
		return false
	}
	return true
}

func (dr DistributionRecord) Add(dr2 DistributionRecord) DistributionRecord {
	dr.Coins = dr.Coins.Add(dr2.Coins...)
	return dr
}

<<<<<<< HEAD
func NewUserClaim(userAddress string, userClaimType DistributionType, time string) UserClaim {
	return UserClaim{UserAddress: userAddress, UserClaimType: userClaimType, Locked: false, UserClaimTime: time}
}

func (uc UserClaim) Validate() bool {
	if len(uc.UserAddress) == 0 {
=======
// The same type is also used Claims
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

func IsValidDistributionType(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "Airdrop":
		return Airdrop, true
	case "LiquidityMining":
		return LiquidityMining, true
	case "ValidatorSubsidy":
		return ValidatorSubsidy, true
	default:
		return 0, false
	}
}

func IsValidClaim(claimType string) (DistributionType, bool) {
	switch claimType {
	case "LiquidityMining":
		return LiquidityMining, true
	case "ValidatorSubsidy":
		return ValidatorSubsidy, true
	default:
		return 0, false
	}
}

// A Distribution object is created for every new distribution type
type Distribution struct {
	DistributionType DistributionType `json:"distribution_type"`
	DistributionName string           `json:"distribution_name"`
	Runner           sdk.AccAddress   `json:"runner"`
}

type Distributions []Distribution

func NewDistribution(t DistributionType, name string, runner sdk.AccAddress) Distribution {
	return Distribution{DistributionType: t, DistributionName: name, Runner: runner}
}

func (d Distribution) Validate() bool {
	if d.DistributionName == "" {
		return false
	}
	_, ok := IsValidDistributionType(d.DistributionType.String())
	if !ok {
		return false
	}
	_, err := sdk.AccAddressFromBech32(d.Runner.String())
	if err != nil {
>>>>>>> develop
		return false
	}

	return true
}

<<<<<<< HEAD
func (uc UserClaim) IsLocked() bool {
	return uc.Locked
}

func NewDistribution(t DistributionType, name string, runner sdk.AccAddress) Distribution {
	return Distribution{DistributionType: t, DistributionName: name, Runner: runner}
=======
func (d Distribution) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s DistributionType :%s`, d.DistributionName, d.DistributionType.String()))
}

type UserClaim struct {
	UserAddress   sdk.AccAddress   `json:"user_address"`
	UserClaimType DistributionType `json:"user_claim_type"`
	UserClaimTime time.Time        `json:"user_claim_time"`
}

type UserClaims []UserClaim

func NewUserClaim(userAddress sdk.AccAddress, userClaimType DistributionType, time time.Time) UserClaim {
	return UserClaim{UserAddress: userAddress, UserClaimType: userClaimType, UserClaimTime: time}
>>>>>>> develop
}

func (d Distribution) Validate() bool {
	if d.DistributionName == "" {
		return false
	}
<<<<<<< HEAD
	_, ok := IsValidDistributionType(d.DistributionType.String())
	if !ok {
		return false
	}
	_, err := sdk.AccAddressFromBech32(d.Runner.String())
	if err != nil {
		return false
	}

=======
	_, err := sdk.AccAddressFromBech32(uc.UserAddress.String())
	if err != nil {
		return false
	}
	_, ok := IsValidClaim(uc.UserClaimType.String())
	if !ok {
		return false
	}
>>>>>>> develop
	return true
}
func GetDistributionStatus(status string) DistributionStatus {
	switch status {
	case "Completed":
		return DistributionStatus_DISTRIBUTION_STATUS_COMPLETED
	case "Pending":
		return DistributionStatus_DISTRIBUTION_STATUS_PENDING
	default:
		return DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED
	}
}

<<<<<<< HEAD
func IsValidClaim(claimType string) (DistributionType, bool) {
	switch claimType {
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}

func IsValidDistribution(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "Airdrop":
		return DistributionType_DISTRIBUTION_TYPE_AIRDROP, true
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
=======
func (uc UserClaim) String() string {
	return strings.TrimSpace(fmt.Sprintf(`UserAddress : %s | UserClaimType : %s`, uc.UserAddress.String(), uc.UserClaimType))
>>>>>>> develop
}
