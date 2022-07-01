//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import (
	"fmt"
)

var (
	RemovalRequestPrefix = []byte{0x08}
	RemovalQueuePrefix   = []byte{0x09}
)

// GetRemovalRequestKey generates a key to store a removal request,
// the key is in the format: lpaddress_id
func GetRemovalRequestKey(request RemovalRequest) []byte {
	key := []byte(fmt.Sprintf("%s_%d", request.Msg.Signer, request.Id))
	return append(RemovalRequestPrefix, key...)
}

func GetRemovalRequestLPPrefix(lpaddress string) []byte {
	key := []byte(fmt.Sprintf("%s", lpaddress))
	return append(RemovalRequestPrefix, key...)
}

func GetRemovalQueueKey(symbol string) []byte {
	key := []byte(fmt.Sprintf("_%s", symbol))
	return append(RemovalRequestPrefix, key...)
}
