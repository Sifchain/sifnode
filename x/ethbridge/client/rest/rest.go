package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	restNetworkDescriptor = "networkDescriptor"
	restBridgeContract    = "bridgeContract"
	restNonce             = "nonce"
	restSymbol            = "symbol"
	restTokenContract     = "tokenContract"
	restEthereumSender    = "ethereumSender"
)

type createEthClaimReq struct {
	BaseReq               rest.BaseReq                  `json:"base_req"`
	NetworkDescriptor     oracletypes.NetworkDescriptor `json:"network_id"`
	BridgeContractAddress string                        `json:"bridge_registry_contract_address"`
	Nonce                 int                           `json:"nonce"`
	Symbol                string                        `json:"symbol"`
	TokenContractAddress  string                        `json:"token_contract_address"`
	EthereumSender        string                        `json:"ethereum_sender"`
	CosmosReceiver        string                        `json:"cosmos_receiver"`
	Validator             string                        `json:"validator"`
	Amount                sdk.Int                       `json:"amount"`
	ClaimType             string                        `json:"claim_type"`
}

type burnOrLockEthReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	NetworkDescriptor string       `json:"network_id"`
	TokenContract     string       `json:"token_contract_address"`
	CosmosSender      string       `json:"cosmos_sender"`
	EthereumReceiver  string       `json:"ethereum_receiver"`
	Amount            sdk.Int      `json:"amount"`
	Symbol            string       `json:"symbol"`
	CethAmount        sdk.Int      `json:"ceth_amount" yaml:"ceth_amount"`
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx client.Context, r *mux.Router, storeName string) {
	getProhechyRoute := fmt.Sprintf(
		"/%s/prophecies/{%s}/{%s}/{%s}/{%s}/{%s}/{%s}",
		storeName, restNetworkDescriptor, restBridgeContract, restNonce,
		restSymbol, restTokenContract, restEthereumSender)

	r.HandleFunc(fmt.Sprintf("/%s/prophecies", storeName), createClaimHandler(cliCtx)).Methods("POST")
	r.HandleFunc(getProhechyRoute, getProphecyHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/burn", storeName), burnOrLockHandler(cliCtx, "burn")).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/lock", storeName), burnOrLockHandler(cliCtx, "lock")).Methods("POST")
}

func createClaimHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createEthClaimReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		bridgeContractAddress := types.NewEthereumAddress(req.BridgeContractAddress)

		tokenContractAddress := types.NewEthereumAddress(req.TokenContractAddress)

		ethereumSender := types.NewEthereumAddress(req.EthereumSender)

		cosmosReceiver, err := sdk.AccAddressFromBech32(req.CosmosReceiver)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		validator, err := sdk.ValAddressFromBech32(req.Validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		claimType, exist := types.ClaimType_value[req.ClaimType]
		if !exist {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrInvalidClaimType.Error())
			return
		}

		ct := types.ClaimType(claimType)

		// create the message
		ethBridgeClaim := types.NewEthBridgeClaim(
			req.NetworkDescriptor, bridgeContractAddress, int64(req.Nonce), req.Symbol,
			tokenContractAddress, ethereumSender, cosmosReceiver, validator, req.Amount, ct)
		msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func getProphecyHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		networkDescriptorString := vars[restNetworkDescriptor]
		networkDescriptor, err := strconv.ParseInt(networkDescriptorString, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		bridgeContract := types.NewEthereumAddress(vars[restBridgeContract])

		nonce := vars[restNonce]
		nonceString, err := strconv.ParseInt(nonce, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		tokenContract := types.NewEthereumAddress(vars[restTokenContract])

		symbol := vars[restSymbol]
		if strings.TrimSpace(symbol) == "" {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "symbol is empty")
			return
		}

		ethereumSender := types.NewEthereumAddress(vars[restEthereumSender])

		bz, err := cliCtx.LegacyAmino.MarshalJSON(
			types.NewQueryEthProphecyRequest(
				oracletypes.NetworkDescriptor(networkDescriptor), bridgeContract, nonceString, symbol, tokenContract, ethereumSender))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryEthProphecy)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func burnOrLockHandler(cliCtx client.Context, lockOrBurn string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnOrLockEthReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		networkDescriptor, err := strconv.Atoi(req.NetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cosmosSender, err := sdk.AccAddressFromBech32(req.CosmosSender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ethereumReceiver := types.NewEthereumAddress(req.EthereumReceiver)

		// create the message
		var msg sdk.Msg
		switch lockOrBurn {
		case "lock":
			msgLock := types.NewMsgLock(oracletypes.NetworkDescriptor(networkDescriptor), cosmosSender, ethereumReceiver, req.Amount, req.Symbol, req.CethAmount)
			msg = &msgLock
		case "burn":
			msgBurn := types.NewMsgBurn(oracletypes.NetworkDescriptor(networkDescriptor), cosmosSender, ethereumReceiver, req.Amount, req.Symbol, req.CethAmount)
			msg = &msgBurn
		}
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
