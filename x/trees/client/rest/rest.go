package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes checks chain id for mainnet safeguard
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	if context.NewCLIContext().ChainID != "sifchain" {
		registerQueryRoutes(cliCtx, r)
		registerTxRoutes(cliCtx, r)
	}
}
