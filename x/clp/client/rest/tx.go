package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	// "bytes"
	// "net/http"

	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

type createPoolReq struct {
	BaseReq             rest.BaseReq `json:"base_req"`
	Signer              string       `json:"signer"`
	ExternalAsset       clp.Asset    `json:"external_asset"`
	NativeAssetAmount   uint         `json:"native_asset_amount"`
	ExternalAssetAmount uint         `json:"external_asset_amount"`
}

func createPooHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPoolReq
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

type addLiquidityReq struct {
	BaseReq             rest.BaseReq `json:"base_req"`
	Signer              string       `json:"signer"`
	ExternalAsset       clp.Asset    `json:"external_asset"`
	NativeAssetAmount   uint         `json:"native_asset_amount"`
	ExternalAssetAmount uint         `json:"external_asset_amount"`
}

func addLiquidityHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addLiquidityReq
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

type removeLiquidityReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Signer        string       `json:"signer"`
	ExternalAsset clp.Asset    `json:"external_asset"`
	WBasisPoints  uint         `json:"w_basis_points"`
	Asymmetry     uint         `json:"asymmetry"`
}

func removeLiquidityHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req removeLiquidityReq
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

type swapReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Signer        string       `json:"signer"`
	SentAsset     clp.Asset    `json:"sent_asset"`
	ReceivedAsset clp.Asset    `json:"received_asset"`
	SentAmount    uint         `json:"sent_amount"`
}

func swapHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req swapReq
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
