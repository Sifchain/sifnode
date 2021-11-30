package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewLegacyQuerier(k KeeperI) sdk.Querier {
	return nil
}

func NewLegacyHandler(k KeeperI) sdk.Handler {
	return nil
}
