package keeper_test
/*
import (
	"github.com/Sifchain/sifnode/simapp"
	"github.com/Sifchain/sifnode/x/clp/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestKeeper_SetPoolIntegration(t *testing.T) {

	pool := test.GenerateRandomPool(1)[0]
	app, ctx := CreateTestApp(false)
	keeper := app.ClpKeeper
	//ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetPool(ctx, pool)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Ticker)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, keeper.ExistsPool(ctx, pool.ExternalAsset.Ticker), true)
}

*/
