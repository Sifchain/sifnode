//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(cliCtx client.Context, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}
