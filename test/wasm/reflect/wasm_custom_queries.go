package reflect

import (
	"encoding/json"
	"strings"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// reflectPlugins needs to be registered in test setup to handle custom query callbacks
func ReflectPlugins() *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: performCustomQuery,
	}
}

func performCustomQuery(_ sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom reflectCustomQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Capitalized != nil {
		msg := strings.ToUpper(custom.Capitalized.Text)
		return json.Marshal(reflectCustomQueryResponse{Msg: msg})
	}
	if custom.Ping != nil {
		return json.Marshal(reflectCustomQueryResponse{Msg: "pong"})
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
}
