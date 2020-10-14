package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	// "bytes"
	// "net/http"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/clp/createPool",
		createPooHandler(cliCtx),
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
		"/clp/swap",
		swapHandler(cliCtx),
	).Methods("POST")
}

type (
	AddLiquidityReq struct {
		BaseReq             rest.BaseReq `json:"base_req"`
		Signer              string       `json:"signer"`
		ExternalAsset       types.Asset  `json:"external_asset"`
		NativeAssetAmount   uint         `json:"native_asset_amount"`
		ExternalAssetAmount uint         `json:"external_asset_amount"`
	}

	RemoveLiquidityReq struct {
		BaseReq       rest.BaseReq `json:"base_req"`
		Signer        string       `json:"signer"`
		ExternalAsset types.Asset  `json:"external_asset"`
		WBasisPoints  uint         `json:"w_basis_points"`
		Asymmetry     uint         `json:"asymmetry"`
	}
	CreatePoolReq struct {
		BaseReq             rest.BaseReq `json:"base_req"`
		Signer              string       `json:"signer"`
		ExternalAsset       types.Asset  `json:"external_asset"`
		NativeAssetAmount   uint         `json:"native_asset_amount"`
		ExternalAssetAmount uint         `json:"external_asset_amount"`
	}
	SwapReq struct {
		BaseReq       rest.BaseReq `json:"base_req"`
		Signer        string       `json:"signer"`
		SentAsset     types.Asset  `json:"sent_asset"`
		ReceivedAsset types.Asset  `json:"received_asset"`
		SentAmount    uint         `json:"sent_amount"`
	}
)

func createPooHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePoolReq
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
		msg := types.NewMsgCreatePool(signer, req.ExternalAsset, req.NativeAssetAmount, req.ExternalAssetAmount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func addLiquidityHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddLiquidityReq
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
		msg := types.NewMsgAddLiquidity(signer, req.ExternalAsset, req.NativeAssetAmount, req.ExternalAssetAmount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func removeLiquidityHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveLiquidityReq
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
		msg := types.NewMsgRemoveLiquidity(signer, req.ExternalAsset, req.WBasisPoints, req.Asymmetry)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func swapHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SwapReq
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
		msg := types.NewMsgSwap(signer, req.SentAsset, req.ReceivedAsset, req.SentAmount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
