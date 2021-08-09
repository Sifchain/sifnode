package rest

import (
	"context"
	"net/http"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/tokenregistry/entries",
		createTokenRegistryEntriesHandler(cliCtx),
	).Methods("GET")
}

type QueryEntriesRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func createTokenRegistryEntriesHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req QueryEntriesRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		queryClient := types.NewQueryClient(cliCtx)
		res, err := queryClient.Entries(context.Background(), &types.QueryEntriesRequest{})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
