package types

import (
	"fmt"
)

const (
	// ModuleName is the name of the module
	ModuleName = "faucet"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName

	FaucetToken      = "rowan"
	MAINNET          = "mainnet"
	TESTNET          = "testnet"
	RequestCoinsType = "request_coins"
	AddCoinsType     = "add_coins"
)

const (
	FaucetPrefix              = "faucet"
	MaxWithdrawAmountPerEpoch = "5000000"
	BlocksPerMinute           = 12
	FaucetResetBlocks         = BlocksPerMinute * 60 * 4 // 4 hours
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetBalanceKey(user string, token string) []byte {
	return []byte(fmt.Sprintf("%s_%s", user, token))
}
