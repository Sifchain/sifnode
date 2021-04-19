package ethbridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// SetCethReceiverAccount
	receiveAccount, err := sdk.AccAddressFromBech32(data.CethReceiveAccount)
	if err != nil && data.CethReceiveAccount != "" {
		panic(err)
	} else {
		keeper.SetCethReceiverAccount(ctx, receiveAccount)
	}

	// AddPeggyTokens
	for _, tokenStr := range data.PeggyTokens.Tokens {
		keeper.AddPeggyToken(ctx, tokenStr)
	}

	return []abci.ValidatorUpdate{}
}


func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	peggyTokens := keeper.GetPeggyToken(ctx)
	receiveAccount := keeper.GetCethReceiverAccount(ctx)

	return &types.GenesisState{
		PeggyTokens: &peggyTokens,
		CethReceiveAccount: receiveAccount.String(),
	}
}

func ValidateGenesis(data types.GenesisState) error {
	return nil
}