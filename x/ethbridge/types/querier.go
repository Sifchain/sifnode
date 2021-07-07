package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQueryEthProphecyRequest creates a new QueryEthProphecyParams
func NewQueryEthProphecyRequest(
	networkDescriptor oracletypes.NetworkDescriptor, bridgeContractAddress EthereumAddress, nonce int64, symbol string,
	tokenContractAddress EthereumAddress, ethereumSender EthereumAddress,
) *QueryEthProphecyRequest {
	return &QueryEthProphecyRequest{
		NetworkDescriptor:             networkDescriptor,
		BridgeContractAddress: bridgeContractAddress.String(),
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContractAddress.String(),
		EthereumSender:        ethereumSender.String(),
	}
}

// NewQueryEthProphecyResponse creates a new QueryEthProphecyResponse instance
func NewQueryEthProphecyResponse(id string, status oracletypes.Status, claims []*EthBridgeClaim) QueryEthProphecyResponse {
	return QueryEthProphecyResponse{
		Id:     id,
		Status: &status,
		Claims: claims,
	}
}
