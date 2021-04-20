package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const (
	// ModuleName is the name of the module
	ModuleName = "dispensation"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute       = ModuleName
	DefaultParamspace  = ModuleName
	MaxRecordsPerBlock = 10
)

var (
	DistributionRecordPrefixPending   = []byte{0x000} // key for storing DistributionRecords pending
	DistributionRecordPrefixCompleted = []byte{0x011} // key for storing DistributionRecords completed
	DistributionsPrefix               = []byte{0x012} // key for storing Distributions [Why does this not work for []byte{0x02} ?!]
)

func GetDistributionRecordKey(status DistributionStatus, name string, recipient string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", name, recipient))
	switch status {
	case Pending:
		return append(DistributionRecordPrefixPending, key...)
	case Completed:
		return append(DistributionRecordPrefixCompleted, key...)
	default:
		return append(DistributionRecordPrefixCompleted, key...)
	}
}
func GetDistributionsKey(name string, distributionType DistributionType) []byte {
	key := []byte(fmt.Sprintf("%s_%d", name, distributionType))
	return append(DistributionsPrefix, key...)
}

func GetDistributionModuleAddress() sdk.AccAddress {
	return supply.NewModuleAddress(ModuleName)
}
