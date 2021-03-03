package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for faucet clients.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown faucet query endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	balance := keeper.GetBankKeeper().GetCoins(ctx, types.GetFaucetModuleAddress())
	res, err := codec.MarshalJSONIndent(keeper.cdc, balance)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
