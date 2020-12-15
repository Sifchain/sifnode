package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

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
