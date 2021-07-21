package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type WhitelistKeeper interface {
	IsDenomWhitelisted(ctx sdk.Context, denom string) (bool, uint)
}
