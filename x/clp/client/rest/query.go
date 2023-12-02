package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/clp/getPool",
		getPoolHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getPools",
		getPoolsHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getLiquidityProvider",
		getLiquidityProviderHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getAssets",
		getAssetsHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getLpList",
		getLpListHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getPmtpParams",
		getPmtpParamsHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getRewardParams",
		getRewardParamsHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/params",
		getParamsHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getRewardsBucket",
		getRewardsBucketHandler(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/clp/getRewardsBuckets",
		getRewardsBucketsHandler(cliCtx),
	).Methods("GET")
}

func getPoolHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		//Generate Router
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool)
		//Generate Params
		var params types.PoolReq
		params.Symbol = r.URL.Query().Get("symbol")

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getLiquidityProviderHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidityProvider)
		var params types.LiquidityProviderReq
		params.Symbol = r.URL.Query().Get("symbol")
		addressString := r.URL.Query().Get("lpAddress")
		lpAddress, err := sdk.AccAddressFromBech32(addressString)
		if err != nil {
			return
		}
		params.LpAddress = lpAddress.String()
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getPoolsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPools)

		var err error
		var limit, offset uint64

		if r.URL.Query().Get("limit") != "" {
			limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if r.URL.Query().Get("offset") != "" {
			offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.PoolsReq{
			Pagination: &query.PageRequest{
				Limit:  limit,
				Offset: offset,
			},
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getAssetsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAssetList)

		var err error
		var limit, offset uint64

		if r.URL.Query().Get("limit") != "" {
			limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if r.URL.Query().Get("offset") != "" {
			offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.AssetListReq{
			Pagination: &query.PageRequest{
				Limit:  limit,
				Offset: offset,
			},
		}

		lpAddress, err := sdk.AccAddressFromBech32(r.URL.Query().Get("lpAddress"))
		if err != nil {
			return
		}

		params.LpAddress = lpAddress.String()

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// http://localhost:1317/clp/getLpList?symbol=catk
func getLpListHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLPList)

		var err error
		var limit, offset uint64
		assetSymbol := r.URL.Query().Get("symbol")

		if r.URL.Query().Get("limit") != "" {
			limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if r.URL.Query().Get("offset") != "" {
			offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.LiquidityProviderListReq{
			Symbol: assetSymbol,
			Pagination: &query.PageRequest{
				Limit:  limit,
				Offset: offset,
			},
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getPmtpParamsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPmtpParams)

		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getRewardParamsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardParams)

		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getParamsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)

		res, height, err := cliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getRewardsBucketHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		//Generate Router
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardsBucket)
		//Generate Params
		var params types.RewardsBucketReq
		params.Denom = r.URL.Query().Get("denom")

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getRewardsBucketsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardsBuckets)

		var err error
		var limit, offset uint64

		if r.URL.Query().Get("limit") != "" {
			limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if r.URL.Query().Get("offset") != "" {
			offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.AllRewardsBucketReq{
			Pagination: &query.PageRequest{
				Limit:  limit,
				Offset: offset,
			},
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
