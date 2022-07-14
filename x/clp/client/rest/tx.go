package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/clp/createPool",
		createPoolHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/clp/addLiquidity",
		addLiquidityHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/clp/removeLiquidity",
		removeLiquidityHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/clp/removeLiquidityUnits",
		removeLiquidityUnitsHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/clp/swap",
		swapHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/clp/decommissionPool",
		decommissionPoolHandler(cliCtx),
	).Methods("POST")
}

type (
	AddLiquidityReq struct {
		BaseReq             rest.BaseReq `json:"base_req"`
		Signer              string       `json:"signer"`                // User who is trying to add liquidity to the pool
		ExternalAsset       types.Asset  `json:"external_asset"`        // ExternalAsset in the pool pair (ex rwn:ceth)
		NativeAssetAmount   sdk.Uint     `json:"native_asset_amount"`   // NativeAssetAmount is the amount of native asset being added
		ExternalAssetAmount sdk.Uint     `json:"external_asset_amount"` // ExternalAssetAmount is the amount of external asset being added
	}

	RemoveLiquidityReq struct {
		BaseReq       rest.BaseReq `json:"base_req"`
		Signer        string       `json:"signer"`         // User who is trying to remove liquidity to the pool
		ExternalAsset types.Asset  `json:"external_asset"` // ExternalAsset in the pool pair (ex rwn:ceth)
		WBasisPoints  sdk.Int      `json:"w_basis_points"` // WBasisPoints determines the amount of asset being withdrawn
		Asymmetry     sdk.Int      `json:"asymmetry"`      // Asymmetry decides the type of asset being withdrawn asymmetry means equal amounts of native and external

	}

	RemoveLiquidityUnitsReq struct {
		BaseReq       rest.BaseReq `json:"base_req"`
		Signer        string       `json:"signer"`         // User who is trying to remove liquidity to the pool
		ExternalAsset types.Asset  `json:"external_asset"` // ExternalAsset in the pool pair (ex rwn:ceth)
		WithdrawUnits sdk.Uint     `json:"withdraw_units"` // WithdrawUnits determines the amount of asset being withdrawn

	}
	CreatePoolReq struct {
		BaseReq             rest.BaseReq `json:"base_req"`
		Signer              string       `json:"signer"`                // User who is trying to create the pool
		ExternalAsset       types.Asset  `json:"external_asset"`        // ExternalAsset in the pool pair (ex rwn:ceth)
		NativeAssetAmount   sdk.Uint     `json:"native_asset_amount"`   // NativeAssetAmount is the amount of native asset being added
		ExternalAssetAmount sdk.Uint     `json:"external_asset_amount"` // ExternalAssetAmount is the amount of external asset being added
	}
	DecommissionPoolReq struct {
		BaseReq rest.BaseReq `json:"base_req"`
		Signer  string       `json:"signer"` // User who is trying to Decommission the pool
		Ticker  string       `json:"ticker"` // ExternalAsset Ticker in the pool pair (ex rwn:ceth ,would be ceth)
	}
	SwapReq struct {
		BaseReq            rest.BaseReq `json:"base_req"`
		Signer             string       `json:"signer"`               // User who is trying to swap
		SentAsset          types.Asset  `json:"sent_asset"`           // Asset which the user is sending ,can be an external asset or RWN
		ReceivedAsset      types.Asset  `json:"received_asset"`       // Asset which the user wants to receive ,can be an external asset or RWN
		SentAmount         sdk.Uint     `json:"sent_amount"`          // Amount of SentAsset being sent
		MinReceivingAmount sdk.Uint     `json:"min_receiving_amount"` // Min amount specified by the user m the swap will not go through if the receiving amount drops below this value
	}
)

//   wallet  < - > abci <-mempool-> tendermint
//   storage > tx
//   /tx hash= []
func createPoolHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePoolReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreatePool(signer, req.ExternalAsset, req.NativeAssetAmount, req.ExternalAssetAmount)

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func decommissionPoolHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DecommissionPoolReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDecommissionPool(signer, req.Ticker)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func addLiquidityHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddLiquidityReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgAddLiquidity(signer, req.ExternalAsset, req.NativeAssetAmount, req.ExternalAssetAmount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func removeLiquidityHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveLiquidityReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRemoveLiquidity(signer, req.ExternalAsset, req.WBasisPoints, req.Asymmetry)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func removeLiquidityUnitsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveLiquidityUnitsReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRemoveLiquidityUnits(signer, req.ExternalAsset, req.WithdrawUnits)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func swapHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SwapReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		signer, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSwap(signer, req.SentAsset, req.ReceivedAsset, req.SentAmount, req.MinReceivingAmount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}
