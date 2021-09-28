package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy           = "prophecies"
	QueryCrosschainFeeConfig   = "crosschainFeeConfig"
	QueryEthereumLockBurnNonce = "ethereumLockBurnNonce"
	QueryWitnessLockBurnNonce  = "witnessLockBurnNonce"
)

// NewQueryEthProphecyRequest creates a new QueryEthProphecyParams
func NewQueryEthProphecyRequest(prophecyID []byte) *QueryEthProphecyRequest {
	return &QueryEthProphecyRequest{
		ProphecyId: prophecyID,
	}
}

// NewQueryEthProphecyResponse creates a new QueryEthProphecyResponse instance
func NewQueryEthProphecyResponse(id []byte, status oracletypes.StatusText, claims []string) QueryEthProphecyResponse {
	return QueryEthProphecyResponse{
		ProphecyId:      id,
		Status:          status,
		ClaimValidators: claims,
	}
}

// NewQueryCrosschainFeeConfigRequest creates a new QueryEthProphecyParams
func NewQueryCrosschainFeeConfigRequest(networkDescriptor oracletypes.NetworkDescriptor) *QueryCrosschainFeeConfigRequest {
	return &QueryCrosschainFeeConfigRequest{
		NetworkDescriptor: networkDescriptor,
	}
}

// NewQueryCrosschainFeeConfigResponse creates a new QueryEthProphecyResponse instance
func NewQueryCrosschainFeeConfigResponse(crosschainFeeConfig oracletypes.CrossChainFeeConfig) QueryCrosschainFeeConfigResponse {
	return QueryCrosschainFeeConfigResponse{
		CrosschainFeeConfig: &crosschainFeeConfig,
	}
}

// NewEthereumLockBurnNonceRequest creates a new QueryLockBurnNonceRequest
func NewEthereumLockBurnNonceRequest(networkDescriptor oracletypes.NetworkDescriptor, relayerValAddress string) *QueryEthereumLockBurnNonceRequest {
	return &QueryEthereumLockBurnNonceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
}

// NewEthereumLockBurnNonceResponse creates a new QueryEthereumLockBurnNonceResponse instance
func NewEthereumLockBurnNonceResponse(lockBurnNonce uint64) QueryEthereumLockBurnNonceResponse {
	return QueryEthereumLockBurnNonceResponse{
		EthereumLockBurnNonce: lockBurnNonce,
	}
}

// NewWitnessLockBurnNonceRequest creates a new QueryWitnessLockBurnNonceRequest
func NewWitnessLockBurnNonceRequest(networkDescriptor oracletypes.NetworkDescriptor, relayerValAddress string) *QueryWitnessLockBurnNonceRequest {
	return &QueryWitnessLockBurnNonceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
}

// NewWitnessLockBurnNonceResponse creates a new QueryWitnessLockBurnNonceResponse instance
func NewWitnessLockBurnNonceResponse(lockBurnNonce uint64) QueryWitnessLockBurnNonceResponse {
	return QueryWitnessLockBurnNonceResponse{
		WitnessLockBurnNonce: lockBurnNonce,
	}
}
