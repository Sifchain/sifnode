package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	FaucetToken = "rowan"
)

const (
	FaucetPrefix              = "faucet"
	MaxWithdrawAmountPerEpoch = "100000"
	BlocksPerMinute           = 12
	FaucetResetBlocks         = BlocksPerMinute * 60 * 4 // 4 hours
)

// Todo : Add MaxWithdrawAmountPerEpoch to Genesis

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetBalanceKey(user sdk.AccAddress, token string) []byte {
	return []byte(fmt.Sprintf("%s_%s", user.String(), token))
}
