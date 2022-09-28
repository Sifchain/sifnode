package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	restProphecyID           = "restProphecyID"
	restNetworkDescriptor    = "networkDescriptor"
	restBridgeContract       = "bridgeContract"
	restSequence             = "sequence"
	restSymbol               = "symbol"
	restTokenContract        = "tokenContract"
	restEthereumSender       = "ethereumSender"
	restRelayerCosmosAddress = "relayerCosmosAddress"
)

type createEthClaimReq struct {
	BaseReq               rest.BaseReq                  `json:"base_req"`
	NetworkDescriptor     oracletypes.NetworkDescriptor `json:"network_descriptor"`
	BridgeContractAddress string                        `json:"bridge_registry_contract_address"`
	Nonce                 int                           `json:"nonce"`
	Symbol                string                        `json:"symbol"`
	TokenContractAddress  string                        `json:"token_contract_address"`
	EthereumSender        string                        `json:"ethereum_sender"`
	CosmosReceiver        string                        `json:"cosmos_receiver"`
	Validator             string                        `json:"validator"`
	Amount                sdk.Int                       `json:"amount"`
	ClaimType             string                        `json:"claim_type"`
	Decimals              int                           `json:"token_decimals"`
	TokenName             string                        `json:"token_name"`
}

type burnOrLockEthReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	NetworkDescriptor string       `json:"network_descriptor"`
	TokenContract     string       `json:"token_contract_address"`
	CosmosSender      string       `json:"cosmos_sender"`
	EthereumReceiver  string       `json:"ethereum_receiver"`
	Amount            sdk.Int      `json:"amount"`
	Symbol            string       `json:"symbol"`
	CrosschainFee     sdk.Int      `json:"cross_chain_fee_amount" yaml:"cross_chain_fee_amount"`
}

type signProphecyReq struct {
	BaseReq               rest.BaseReq `json:"base_req"`
	NetworkDescriptor     string       `json:"network_descriptor"`
	CosmosSender          string       `json:"cosmos_sender"`
	EthereumSignerAddress string       `json:"ethereum_signer_address"`
	Signature             string       `json:"signature"`
	ProphecyID            string       `json:"prophecy_id"`
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx client.Context, r *mux.Router, storeName string) {
	getProhechyRoute := fmt.Sprintf(
		"/%s/prophecies/{%s}/{%s}/{%s}/{%s}/{%s}/{%s}",
		storeName, restNetworkDescriptor, restBridgeContract, restSequence,
		restSymbol, restTokenContract, restEthereumSender)

	getCrosschainFeeConfigRoute := fmt.Sprintf("/%s/crosschainFeeConfig/{%s}", storeName, restNetworkDescriptor)
	getEthereumLockBurnNonceRoute := fmt.Sprintf("/%s/ethereumLockBurnNonce/{%s}/{%s}", storeName, restNetworkDescriptor, restRelayerCosmosAddress)
	getWitnessLockBurnNonceRoute := fmt.Sprintf("/%s/witnessLockBurnNonce/{%s}/{%s}", storeName, restNetworkDescriptor, restRelayerCosmosAddress)
	getGlobalSequenceBlockNumberRoute := fmt.Sprintf("/%s/globalSequenceBlockNumber/{%s}/{%s}", storeName, restNetworkDescriptor, restSequence)
	getProphciesCompletedRoute := fmt.Sprintf("/%s/ProphciesCompleted/{%s}/{%s}", storeName, restNetworkDescriptor, restSequence)

	r.HandleFunc(fmt.Sprintf("/%s/prophecies", storeName), createClaimHandler(cliCtx)).Methods("POST")
	r.HandleFunc(getProhechyRoute, getProphecyHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(getCrosschainFeeConfigRoute, getCrosschainFeeConfigHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/burn", storeName), burnOrLockHandler(cliCtx, "burn")).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/lock", storeName), burnOrLockHandler(cliCtx, "lock")).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/signProphecy", storeName), signProphecyHandler(cliCtx)).Methods("POST")
	r.HandleFunc(getEthereumLockBurnNonceRoute, getEthereumLockBurnSequenceHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(getWitnessLockBurnNonceRoute, getWitnessLockBurnSequenceHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(getGlobalSequenceBlockNumberRoute, getQueryGlobalSequenceBlockNumberHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(getProphciesCompletedRoute, getProphciesCompletedHandler(cliCtx, storeName)).Methods("GET")
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
			req.NetworkDescriptor, bridgeContractAddress, uint64(req.Nonce), req.Symbol,
			tokenContractAddress, ethereumSender, cosmosReceiver, validator, req.Amount, ct, req.TokenName, uint8(req.Decimals))
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

		restProphecyID := []byte(vars[restProphecyID])

		bz, err := cliCtx.LegacyAmino.MarshalJSON(types.NewQueryEthProphecyRequest(restProphecyID))
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

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(req.NetworkDescriptor)
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
			msgLock := types.NewMsgLock(networkDescriptor, cosmosSender, ethereumReceiver, req.Amount, req.Symbol, req.CrosschainFee)
			msg = &msgLock
		case "burn":
			msgBurn := types.NewMsgBurn(networkDescriptor, cosmosSender, ethereumReceiver, req.Amount, req.Symbol, req.CrosschainFee)
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

func signProphecyHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signProphecyReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(req.NetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = sdk.AccAddressFromBech32(req.CosmosSender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		prophecyID := req.ProphecyID
		ethereumAddress := req.EthereumSignerAddress
		signature := req.Signature

		// create the message
		msg := types.NewMsgSignProphecy(req.CosmosSender, networkDescriptor, []byte(prophecyID), ethereumAddress, signature)

		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func getCrosschainFeeConfigHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		restNetworkDescriptor := vars[restNetworkDescriptor]

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(restNetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(types.NewQueryCrosschainFeeConfigRequest(networkDescriptor))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryCrosschainFeeConfig)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getEthereumLockBurnSequenceHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		restNetworkDescriptor := vars[restNetworkDescriptor]

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(restNetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		valAddress := vars[restRelayerCosmosAddress]

		bz, err := cliCtx.LegacyAmino.MarshalJSON(types.NewEthereumLockBurnSequenceRequest(networkDescriptor, valAddress))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryEthereumLockBurnSequence)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getWitnessLockBurnSequenceHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		restNetworkDescriptor := vars[restNetworkDescriptor]

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(restNetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		valAddress := vars[restRelayerCosmosAddress]

		bz, err := cliCtx.LegacyAmino.MarshalJSON(types.NewWitnessLockBurnSequenceRequest(networkDescriptor, valAddress))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryWitnessLockBurnSequence)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getQueryGlobalSequenceBlockNumberHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		restNetworkDescriptor := vars[restNetworkDescriptor]

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(restNetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		globalSequence, err := strconv.ParseInt(restNetworkDescriptor, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(types.NewQueryGlobalSequenceBlockNumberRequest(networkDescriptor, uint64(globalSequence)))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryGlobalSequenceBlockNumber)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getProphciesCompletedHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		restNetworkDescriptor := vars[restNetworkDescriptor]

		networkDescriptor, err := oracletypes.ParseNetworkDescriptor(restNetworkDescriptor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		globalSequence, err := strconv.ParseInt(restSequence, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(types.NewPropheciesCompletedRequest(networkDescriptor, uint64(globalSequence)))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryPropheciesCompleted)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
