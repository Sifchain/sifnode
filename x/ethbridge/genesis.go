package ethbridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{
		CethReceiveAccount: "",
		PeggyTokens:        []string{},
		Blacklist:          []string{},
		Pause:              &types.Pause{IsPaused: false},
	}
}

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// SetCethReceiverAccount
	if data.CethReceiveAccount != "" {
		receiveAccount, err := sdk.AccAddressFromBech32(data.CethReceiveAccount)
		if err != nil {
			panic(err)
		}
		keeper.SetCethReceiverAccount(ctx, receiveAccount)
	}

	// AddPeggyTokens
	if data.PeggyTokens != nil {
		for _, tokenStr := range data.PeggyTokens {
			keeper.AddPeggyToken(ctx, tokenStr)
		}
	}

	// Set blacklisted addresses
	for _, address := range data.Blacklist {
		keeper.SetBlacklistAddress(ctx, address)
	}

	// Set pause
	if data.Pause != nil {
		keeper.SetPause(ctx, data.Pause)
	} else {
		keeper.SetPause(ctx, &types.Pause{IsPaused: false})
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	peggyTokens := keeper.GetPeggyToken(ctx)
	receiveAccount := keeper.GetCethReceiverAccount(ctx)
	blacklist := keeper.GetBlacklist(ctx)
	isPaused := keeper.IsPaused(ctx)

	// create pause
	pause := types.Pause{IsPaused: isPaused}

	return &types.GenesisState{
		PeggyTokens:        peggyTokens.Tokens,
		CethReceiveAccount: receiveAccount.String(),
		Blacklist:          blacklist,
		Pause:              &pause,
	}
}

func ValidateGenesis(data types.GenesisState) error {
	return nil
}
