package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	// "bytes"
	// "net/http"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	// "github.com/Sifchain/sifnode/x/faucet/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/faucet/leak", createLeakHandler(cliCtx)).Methods("POST")
}

type createLeakRequest struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Minter  sdk.AccAddress `json:"minter" yaml:"minter"`
	Amount  int            `json:"amount" yaml:"amount"`
}

func createLeakHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
