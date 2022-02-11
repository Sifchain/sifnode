package wasm

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Plugins needs to be registered to handle custom query callbacks
func Plugins(clpKeeper clpkeeper.Keeper) *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: SifchainQuerier(clpKeeper),
	}
}

func SifchainQuerier(clpKeeper clpkeeper.Keeper) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(context sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom SifchainQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		switch {
		case custom.Pool != nil:
			return PerformPoolQuery(clpKeeper, context, custom.Pool)
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Sifchain query variant")
	}
}

func PerformPoolQuery(
	clpKeeper clpkeeper.Keeper,
	ctx sdk.Context,
	poolQuery *PoolQuery) ([]byte, error) {

	pool, err := clpKeeper.GetPool(ctx, poolQuery.ExternalAsset)
	if err != nil {
		return nil, err
	}

	resp := PoolResponse{
		ExternalAsset:        pool.ExternalAsset.Symbol,
		ExternalAssetBalance: pool.ExternalAssetBalance.String(),
		NativeAssetBalance:   pool.NativeAssetBalance.String(),
		PoolUnits:            pool.PoolUnits.String(),
	}

	return json.Marshal(resp)
}
