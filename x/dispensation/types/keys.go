package types

import "fmt"

const (
	// ModuleName is the name of the module
	ModuleName = "dispensation"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	DistributionRecordPrefix = []byte{0x00} // key for storing DistributionRecords
	AirdropRecordPrefix      = []byte{0x01} // key for storing airdropRecords
)

func GetDistributionRecordKey(airdropName string, recipient string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", airdropName, recipient))
	return append(DistributionRecordPrefix, key...)
}
func GetAirdropRecordKey(airdropName string) []byte {
	key := []byte(fmt.Sprintf("%s", airdropName))
	return append(AirdropRecordPrefix, key...)
}
