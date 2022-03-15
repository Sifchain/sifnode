package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tenderminttypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestEndBlock(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	// Setup reward period
	params := app.ClpKeeper.GetParams(ctx)
	allocation := sdk.NewUint(10)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	err := app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:        &types.Asset{Symbol: "atom"},
		NativeAssetBalance:   sdk.NewUint(1000),
		ExternalAssetBalance: sdk.NewUint(1000),
		PoolUnits:            sdk.NewUint(1000),
	})
	require.NoError(t, err)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.False(t, startingSupply.Equal(periodOneSupply))
	require.True(t, periodOneSupply.Sub(startingSupply).Sub(sdk.NewCoin("rowan", sdk.NewInt(10))).IsZero())
	// check pool has expected increase
	pool, err := app.ClpKeeper.GetPool(ctx, "atom")
	require.NoError(t, err)
	require.Equal(t, "1010", pool.NativeAssetBalance.String())
	// continue through a non reward period
	for block := 11; block <= 20; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply is unchanged
	supplyCheck := app.BankKeeper.GetSupply(ctx, "rowan")
	require.True(t, supplyCheck.Equal(periodOneSupply))
}
