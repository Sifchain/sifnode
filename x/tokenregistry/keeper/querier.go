package keeper

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewLegacyQuerier(keeper types.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "Token Registry Legacy Querier No Longer Available")
	}
}
