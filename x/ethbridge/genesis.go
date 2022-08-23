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

	for key, value := range data.EthereumLockBurnSequence {
		keeper.SetSequenceViaRawKey(ctx, []byte(key), value)
	}

	for key, value := range data.GlobalNonce {
		keeper.SetGlobalSequenceViaRawKey(ctx, key, value)
	}

	for key, value := range data.GlobalNonceBlockNumber {
		keeper.SetGlobalSequenceToBlockNumberViaRawKey(ctx, key, value)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		CrosschainFeeReceiveAccount: keeper.GetCrossChainFeeReceiverAccount(ctx).String(),
		EthereumLockBurnSequence:    keeper.GetEthereumLockBurnSequences(ctx),
		GlobalNonce:                 keeper.GetGlobalSequences(ctx),
		GlobalNonceBlockNumber:      keeper.GetGlobalSequenceToBlockNumbers(ctx),
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
