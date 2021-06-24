package rest

import (
	"fmt"
<<<<<<< HEAD
	"net/http"

=======
>>>>>>> develop
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
<<<<<<< HEAD
=======
	"net/http"
>>>>>>> develop
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Input :No parameters
	// Output : List of all distributions created on the network
	r.HandleFunc(
		"/dispensation/getDistributions",
		getDistributionsHandler(cliCtx),
	).Methods("GET")
	// Input : User Address (address)
	// Output : List of all Distribution Records for the address (Completed and Pending)
	r.HandleFunc(
		"/dispensation/getDrForAddress",
		getDrForRecipientHandler(cliCtx),
	).Methods("GET")
	// Input : Distribution Name (distName)
	//		   Status (status)	 [accepts Pending and Completed]
	// Output : List of all Distribution Records Distribution name and status provided
	r.HandleFunc(
		"/dispensation/getDrForDistName",
		getDrForDistHandler(cliCtx),
	).Methods("GET")
	// Input : Claim type (type) [Accepts ValidatorSubsidy and LiquidityMining]
	// Output : List of all Claims for the type.
	r.HandleFunc(
		"/dispensation/getClaims",
		getClaimsHandler(cliCtx),
	).Methods("GET")
}

func getDistributionsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		//Generate Router
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllDistributions)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getDrForRecipientHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		//Generate Router
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordsByRecipient)
		//Generate Params
		var params types.QueryRecordsByRecipientAddr
		addressString := r.URL.Query().Get("address")
		address, err := sdk.AccAddressFromBech32(addressString)
		if err != nil {
			return
		}
		params.Address = address
		bz, err := cliCtx.Codec.MarshalJSON(params)
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

func getDrForDistHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		//Generate Router
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordsByDistrName)
		//Generate Params
		var params types.QueryRecordsByDistributionName
		params.DistributionName = r.URL.Query().Get("distName")
		status, ok := types.IsValidStatus(r.URL.Query().Get("status"))
		if !ok {
			return
		}
		params.Status = status
		bz, err := cliCtx.Codec.MarshalJSON(params)
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

func getClaimsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		//Generate Router
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryClaimsByType)
		//Generate Params
		var params types.QueryUserClaims
		claimType, ok := types.IsValidClaim(r.URL.Query().Get("type"))
		if !ok {
			return
		}
		params.UserClaimType = claimType
		bz, err := cliCtx.Codec.MarshalJSON(params)
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
