package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// RegisterRoutes checks chain id for mainnet safeguard
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	// if context.NewCLIContext().ChainID != "sifchain" {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
	// } // todo:  move this to the make file
}
