package rest

import (
	"fmt"
	"net/http"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/tokenregistry/entries",
		getTokenRegistryHandler(cliCtx),
	).Methods("GET")
}

type QueryEntriesRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func getTokenRegistryHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEntries)

		bz, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var res types.QueryEntriesResponse
		err = types.ModuleCdc.UnmarshalJSON(bz, &res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
