package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Sifchain/sifnode/x/faucet/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/faucet/requestCoins", createCoinsHandler(cliCtx)).Methods("POST")
}

type (
	CoinsReq struct {
		BaseReq rest.BaseReq `json:"base_req"`
		Signer  string       `json:"signer"`
		Amount  string       `json:"amount" yaml:"amount"`
	}
)

func createCoinsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		coins, err := sdk.ParseCoins(req.Amount)
		msg := types.NewMsgRequestCoins(signer, coins)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
