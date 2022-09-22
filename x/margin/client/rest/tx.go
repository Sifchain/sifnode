package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/margin/open",
		openHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/margin/close",
		closeHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/margin/forceClose",
		forceCloseHandler(cliCtx),
	).Methods("POST")
}

type (
	OpenReq struct {
		BaseReq          rest.BaseReq   `json:"base_req"`
		Signer           string         `json:"signer"`            // User who is trying to open margin position
		CollateralAsset  clptypes.Asset `json:"collateral_asset"`  // CollateralAsset for margin position
		CollateralAmount sdk.Uint       `json:"collateral_amount"` // CollateralAmount is the amount of collateral being added
		BorrowAsset      clptypes.Asset `json:"borrow_asset"`      // BorrowAsset is asset being borrowed in margin position
		Position         types.Position `json:"position"`          // Position type for margin position
	}

	CloseReq struct {
		BaseReq rest.BaseReq `json:"base_req"`
		Signer  string       `json:"signer"` // User who is trying to close margin position
		// nolint:revive
		Id uint64 `json:"id"` // Id of the mtp

	}

	ForceCloseReq struct {
		BaseReq    rest.BaseReq `json:"base_req"`
		Signer     string       `json:"signer"`      // User who is trying to close margin position
		MtpAddress string       `json:"mtp_address"` // MtpAddress for position to force close
		// nolint:revive
		Id uint64 `json:"id"` // Id of the mtp

	}
)

func openHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OpenReq
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

		msg := types.MsgOpen{
			Signer:           signer.String(),
			CollateralAsset:  req.CollateralAsset.String(),
			CollateralAmount: req.CollateralAmount,
			BorrowAsset:      req.BorrowAsset.String(),
			Position:         req.Position,
		}

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func closeHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CloseReq
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

		msg := types.MsgClose{
			Signer: signer.String(),
			Id:     req.Id,
		}

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func forceCloseHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ForceCloseReq
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

		mtpAddress, err := sdk.AccAddressFromBech32(req.MtpAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.MsgForceClose{
			Signer:     signer.String(),
			MtpAddress: mtpAddress.String(),
			Id:         req.Id,
		}

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}
