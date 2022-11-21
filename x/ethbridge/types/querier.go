package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy               = "prophecies"
	QueryCrosschainFeeConfig       = "crosschainFeeConfig"
	QueryEthereumLockBurnSequence  = "ethereumLockBurnSequence"
	QueryWitnessLockBurnSequence   = "witnessLockBurnSequence"
	QueryGlobalSequenceBlockNumber = "globalSequenceBlockNumber"
	QueryPropheciesCompleted       = "propheciesCompleted"

	QueryBlacklist = "blacklist"
	QueryPause       = "pause"
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

// NewEthereumLockBurnSequenceRequest creates a new QueryLockBurnSequenceRequest
func NewEthereumLockBurnSequenceRequest(networkDescriptor oracletypes.NetworkDescriptor, relayerValAddress string) *QueryEthereumLockBurnSequenceRequest {
	return &QueryEthereumLockBurnSequenceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
}

// NewEthereumLockBurnSequenceResponse creates a new QueryEthereumLockBurnSequenceResponse instance
func NewEthereumLockBurnSequenceResponse(LockBurnSequence uint64) QueryEthereumLockBurnSequenceResponse {
	return QueryEthereumLockBurnSequenceResponse{
		EthereumLockBurnSequence: LockBurnSequence,
	}
}

// NewWitnessLockBurnSequenceRequest creates a new QueryWitnessLockBurnSequenceRequest
func NewWitnessLockBurnSequenceRequest(networkDescriptor oracletypes.NetworkDescriptor, relayerValAddress string) *QueryWitnessLockBurnSequenceRequest {
	return &QueryWitnessLockBurnSequenceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
}

// NewWitnessLockBurnSequenceResponse creates a new QueryWitnessLockBurnSequenceResponse instance
func NewWitnessLockBurnSequenceResponse(LockBurnSequence uint64) QueryWitnessLockBurnSequenceResponse {
	return QueryWitnessLockBurnSequenceResponse{
		WitnessLockBurnSequence: LockBurnSequence,
	}
}

// NewQueryGlobalSequenceBlockNumberRequest creates a new QueryGlobalSequenceBlockNumberRequest
func NewQueryGlobalSequenceBlockNumberRequest(networkDescriptor oracletypes.NetworkDescriptor, globalSequence uint64) *QueryGlobalSequenceBlockNumberRequest {
	return &QueryGlobalSequenceBlockNumberRequest{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    globalSequence,
	}
}

// NewGlobalSequenceBlockNumberResponse creates a new QueryWitnessLockBurnSequenceResponse instance
func NewGlobalSequenceBlockNumberResponse(blockNumber uint64) QueryGlobalSequenceBlockNumberResponse {
	return QueryGlobalSequenceBlockNumberResponse{
		BlockNumber: blockNumber,
	}
}

// NewPropheciesCompletedRequest creates a new NewGlobalSequenceBlockNumberResponse instance
func NewPropheciesCompletedRequest(networkDescriptor oracletypes.NetworkDescriptor, globalSequence uint64) *QueryPropheciesCompletedRequest {
	return &QueryPropheciesCompletedRequest{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    globalSequence,
	}
}

// NewQueryPropheciesCompletedResponse creates a new QueryWitnessLockBurnSequenceResponse instance
func NewQueryPropheciesCompletedResponse(prophecyInfo []*oracletypes.ProphecyInfo) QueryPropheciesCompletedResponse {
	return QueryPropheciesCompletedResponse{
		ProphecyInfo: prophecyInfo,
	}
}
