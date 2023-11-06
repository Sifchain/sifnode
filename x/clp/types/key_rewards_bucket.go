package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RewardsBucketKeyPrefix is the prefix to retrieve all RewardsBucket
	RewardsBucketKeyPrefix = "RewardsBucket/value/"
)

// RewardsBucketKey returns the store key to retrieve a RewardsBucket from the index fields
func RewardsBucketKey(
	denom string,
) []byte {
	var key []byte

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}
