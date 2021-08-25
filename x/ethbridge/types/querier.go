package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy         = "prophecies"
	QueryCrosschainFeeConfig = "crosschainFeeConfig"
	QueryLockBurnNonce       = "lockBurnNonce"
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

// NewLockBurnNonceRequest creates a new QueryLockBurnNonceRequest
func NewLockBurnNonceRequest(networkDescriptor oracletypes.NetworkDescriptor, relayerValAddress string) *QueryLockBurnNonceRequest {
	return &QueryLockBurnNonceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
}

// NewLockBurnNonceResponse creates a new QueryLockBurnNonceResponse instance
func NewLockBurnNonceResponse(lockBurnNonce uint64) QueryLockBurnNonceResponse {
	return QueryLockBurnNonceResponse{
		LockBurnNonce: lockBurnNonce,
	}
}
