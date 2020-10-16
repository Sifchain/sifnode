package keeper

import (
  // this line is used by starport scaffolding
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/Sifchain/sifnode/x/sifnode/types"
)

// NewQuerier creates a new querier for sifnode clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
    // this line is used by starport scaffolding # 2
		case types.QueryListUser:
			return listUser(ctx, k)
		case types.QueryGetUser:
			return getUser(ctx, path[1:], k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown sifnode query endpoint")
		}
	}
}