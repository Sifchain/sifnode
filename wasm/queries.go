package wasm

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Plugins needs to be registered to handle custom query callbacks
func Plugins() *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: PerformQuery,
	}
}

func PerformQuery(_ sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom SifchainQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Ping != nil {
		return json.Marshal(SifchainQueryResponse{Msg: "pong"})
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Sifchain query variant")
}
