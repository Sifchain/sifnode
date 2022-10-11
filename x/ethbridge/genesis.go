package ethbridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// SetCrossChainFeeReceiverAccount
	if data.CrosschainFeeReceiveAccount != "" {
		receiveAccount, err := sdk.AccAddressFromBech32(data.CrosschainFeeReceiveAccount)
		if err != nil {
			panic(err)
		}
		keeper.SetCrossChainFeeReceiverAccount(ctx, receiveAccount)
	}

	for _, value := range data.GenesisEthereumLockBurnSequence {
		keeper.SetEthereumLockBurnSequence(ctx, value.EthereumLockBurnSequenceKey.NetworkDescriptor,
			value.EthereumLockBurnSequenceKey.ValidatorAddress,
			value.EthereumLockBurnSequence.EthereumLockBurnSequence)
	}

	for _, value := range data.GenesisGlobalSequence {
		keeper.UpdateGlobalSequence(ctx, value.NetworkDescriptor, value.GlobalSequence.GlobalSequence)
	}

	for _, value := range data.GlobalNonceBlockNumber {
		keeper.SetGlobalSequenceToBlockNumber(ctx, value.GlobalSequenceKey.NetworkDescriptor, value.GlobalSequenceKey.GlobalSequence, value.BlockNumber.BlockNumber)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		CrosschainFeeReceiveAccount:     keeper.GetCrossChainFeeReceiverAccount(ctx).String(),
		GenesisEthereumLockBurnSequence: keeper.GetEthereumLockBurnSequences(ctx),
		GenesisGlobalSequence:           keeper.GetGlobalSequences(ctx),
		GlobalNonceBlockNumber:          keeper.GetGlobalSequenceToBlockNumbers(ctx),
	}
}

// ValidateGenesis check all values in genesis are valid
func ValidateGenesis(data types.GenesisState) error {
	if data.CrosschainFeeReceiveAccount == "" {
		return nil
	}
	_, err := sdk.AccAddressFromBech32(data.CrosschainFeeReceiveAccount)
	return err
}
