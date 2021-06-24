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
<<<<<<< HEAD
	QuerierRoute       = ModuleName
	DefaultParamspace  = ModuleName
	MaxRecordsPerBlock = 10
)

var (
	DistributionRecordPrefixPending   = []byte{0x000} // key for storing DistributionRecords pending
	DistributionRecordPrefixCompleted = []byte{0x011} // key for storing DistributionRecords completed
	DistributionsPrefix               = []byte{0x01}  // key for storing Distributions
	UserClaimPrefix                   = []byte{0x02}  // key for storing user claims
)

func GetDistributionRecordKey(status DistributionStatus, name string, recipient string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", name, recipient))
	switch status {
	case DistributionStatus_DISTRIBUTION_STATUS_PENDING:
		return append(DistributionRecordPrefixPending, key...)
	case DistributionStatus_DISTRIBUTION_STATUS_COMPLETED:
		return append(DistributionRecordPrefixCompleted, key...)
	default:
		return append(DistributionRecordPrefixCompleted, key...)
	}
=======
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
	TokenSupported    = "rowan"
)

var (
	DistributionRecordPrefix = []byte{0x00} // key for storing DistributionRecords
	DistributionsPrefix      = []byte{0x01} // key for storing airdropRecords
	// Would need to verify usage for byte 02,03 and 04
	UserClaimPrefix                = []byte{0x05} // key for storing airdropRecords
	DistributionRecordPrefixFailed = []byte{0x06} // key for storing DistributionRecords
)

// A distribution records is unique for name_recipientAddress
// key format  :  Height_Distributor_Receipient_DistributionType
func GetDistributionRecordKey(name string, recipient string, distributionType string) []byte {
	key := []byte(fmt.Sprintf("%s_%s_%s", name, distributionType, recipient))
	return append(DistributionRecordPrefix, key...)
>>>>>>> develop
}

// A distribution faile records is the similar to GetDistributionRecordKey , but uses a different prefix
func GetDistributionRecordFailedKey(name string, recipient string, distributionType string) []byte {
	key := []byte(fmt.Sprintf("%s_%s_%s", name, distributionType, recipient))
	return append(DistributionRecordPrefixFailed, key...)
}

// A distribution  is unique for name_distributionType
func GetDistributionsKey(name string, distributionType DistributionType) []byte {
	key := []byte(fmt.Sprintf("%s_%d", name, distributionType))
	return append(DistributionsPrefix, key...)
}

// A claim is unique for userAddress_userClaimType
func GetUserClaimKey(userAddress string, userClaimType DistributionType) []byte {
	key := []byte(fmt.Sprintf("%s_%d", userAddress, userClaimType))
	return append(UserClaimPrefix, key...)
}

func GetDistributionModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}
