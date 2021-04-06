package types

import (
	"encoding/json"
	"fmt"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQueryEthProphecyRequest creates a new QueryEthProphecyParams
func NewQueryEthProphecyRequest(
	ethereumChainID int64, bridgeContractAddress EthereumAddress, nonce int64, symbol string,
	tokenContractAddress EthereumAddress, ethereumSender EthereumAddress,
) *QueryEthProphecyRequest {
	return &QueryEthProphecyRequest{
		EthereumChainId:       ethereumChainID,
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

// String implements fmt.Stringer interface
func (response QueryEthProphecyResponse) String() string {
	prophecyJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}
