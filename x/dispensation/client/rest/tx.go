package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/dispensation/airdrop",
		airdropHandler(cliCtx),
	).Methods("POST")
}

func airdropHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
