package keeper

import (
	"github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewLegacyQuerier(keeper types.Keeper) sdk.Querier {
	querier := Querier{keeper}
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryEntries:
			return queryDenoms(ctx, querier)

		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown dispensation query endpoint")
		}
	}
}
func queryDenoms(ctx sdk.Context, querier Querier) ([]byte, error) {
	res, err := querier.Entries(sdk.WrapSDKContext(ctx), &types.QueryEntriesRequest{})
	if err != nil {
		return nil, err
	}

	return types.ModuleCdc.MarshalJSON(res)
}
