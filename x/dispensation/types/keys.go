package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName is the name of the module
	ModuleName                = "dispensation"
	MsgTypeCreateUserClaim    = "createUserClaim"
	MsgTypeCreateDistribution = "createDistribution"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute       = ModuleName
	DefaultParamspace  = ModuleName
	MaxRecordsPerBlock = 10
	TokenSupported     = "rowan"
)

var (
	DistributionRecordPrefixPending   = []byte{0x000} // key for storing DistributionRecords pending
	DistributionRecordPrefixCompleted = []byte{0x011} // key for storing DistributionRecords completed
	DistributionRecordPrefixFailed    = []byte{0x012} // key for storing DistributionRecords failed
	DistributionsPrefix               = []byte{0x01}  // key for storing Distributions
	UserClaimPrefix                   = []byte{0x02}  // key for storing user claims
)

func GetDistributionRecordKey(status DistributionStatus, name string, recipient string, distributionType DistributionType) []byte {
	key := []byte(fmt.Sprintf("%s_%d_%s", name, distributionType, recipient))
	switch status {
	case DistributionStatus_DISTRIBUTION_STATUS_PENDING:
		return append(DistributionRecordPrefixPending, key...)
	case DistributionStatus_DISTRIBUTION_STATUS_COMPLETED:
		return append(DistributionRecordPrefixCompleted, key...)
	case DistributionStatus_DISTRIBUTION_STATUS_FAILED:
		return append(DistributionRecordPrefixFailed, key...)
	default:
		return append(DistributionRecordPrefixCompleted, key...)
	}
}
func GetDistributionsKey(name string, distributionType DistributionType) []byte {
	key := []byte(fmt.Sprintf("%s_%d", name, distributionType))
	return append(DistributionsPrefix, key...)
}

func GetUserClaimKey(userAddress string, userClaimType DistributionType) []byte {
	key := []byte(fmt.Sprintf("%s_%d", userAddress, userClaimType))
	return append(UserClaimPrefix, key...)
}

func GetDistributionModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}
