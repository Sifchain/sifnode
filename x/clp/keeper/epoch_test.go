package keeper_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestKeeper_SetPmtpEpoch(t *testing.T) {
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	const poolAsset = "eth"
	nativeBalance := sdk.NewInt(10000)
	externalBalance := sdk.NewInt(10000)

	ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		balances := []banktypes.Balance{
			{
				Address: address,
				Coins: sdk.Coins{
					sdk.NewCoin(poolAsset, externalBalance),
					sdk.NewCoin("rowan", nativeBalance),
				},
			},
		}
		bankGs := banktypes.DefaultGenesisState()
		bankGs.Balances = append(bankGs.Balances, balances...)
		bz, _ := app.AppCodec().MarshalJSON(bankGs)
		genesisState["bank"] = bz

		return genesisState
	})

	params := types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 1000,
	}

	app.ClpKeeper.SetPmtpEpoch(ctx, params)

	got := app.ClpKeeper.GetPmtpEpoch(ctx)

	require.Equal(t, got, types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 1000,
	})
}

func TestKeeper_DecrementEpochCounter(t *testing.T) {
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	const poolAsset = "eth"
	nativeBalance := sdk.NewInt(10000)
	externalBalance := sdk.NewInt(10000)

	ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		balances := []banktypes.Balance{
			{
				Address: address,
				Coins: sdk.Coins{
					sdk.NewCoin(poolAsset, externalBalance),
					sdk.NewCoin("rowan", nativeBalance),
				},
			},
		}
		bankGs := banktypes.DefaultGenesisState()
		bankGs.Balances = append(bankGs.Balances, balances...)
		bz, _ := app.AppCodec().MarshalJSON(bankGs)
		genesisState["bank"] = bz

		return genesisState
	})

	params := types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 1000,
	}

	app.ClpKeeper.SetPmtpEpoch(ctx, params)

	app.ClpKeeper.DecrementEpochCounter(ctx)

	got := app.ClpKeeper.GetPmtpEpoch(ctx)

	require.Equal(t, got, types.PmtpEpoch{
		EpochCounter: 999,
		BlockCounter: 1000,
	})
}

func TestKeeper_DecrementBlockCounter(t *testing.T) {
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	const poolAsset = "eth"
	nativeBalance := sdk.NewInt(10000)
	externalBalance := sdk.NewInt(10000)

	ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		balances := []banktypes.Balance{
			{
				Address: address,
				Coins: sdk.Coins{
					sdk.NewCoin(poolAsset, externalBalance),
					sdk.NewCoin("rowan", nativeBalance),
				},
			},
		}
		bankGs := banktypes.DefaultGenesisState()
		bankGs.Balances = append(bankGs.Balances, balances...)
		bz, _ := app.AppCodec().MarshalJSON(bankGs)
		genesisState["bank"] = bz

		return genesisState
	})

	params := types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 1000,
	}

	app.ClpKeeper.SetPmtpEpoch(ctx, params)

	app.ClpKeeper.DecrementBlockCounter(ctx)

	got := app.ClpKeeper.GetPmtpEpoch(ctx)

	require.Equal(t, got, types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 999,
	})
}

func TestKeeper_SetBlockCounter(t *testing.T) {
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	const poolAsset = "eth"
	nativeBalance := sdk.NewInt(10000)
	externalBalance := sdk.NewInt(10000)

	ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		balances := []banktypes.Balance{
			{
				Address: address,
				Coins: sdk.Coins{
					sdk.NewCoin(poolAsset, externalBalance),
					sdk.NewCoin("rowan", nativeBalance),
				},
			},
		}
		bankGs := banktypes.DefaultGenesisState()
		bankGs.Balances = append(bankGs.Balances, balances...)
		bz, _ := app.AppCodec().MarshalJSON(bankGs)
		genesisState["bank"] = bz

		return genesisState
	})

	params := types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 1000,
	}

	app.ClpKeeper.SetPmtpEpoch(ctx, params)

	app.ClpKeeper.SetBlockCounter(ctx, 2000)

	got := app.ClpKeeper.GetPmtpEpoch(ctx)

	require.Equal(t, got, types.PmtpEpoch{
		EpochCounter: 1000,
		BlockCounter: 2000,
	})
}
