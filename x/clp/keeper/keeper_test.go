package keeper_test

import (
	"errors"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestKeeper_Errors(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	num := clpKeeper.GetMinCreatePoolThreshold(ctx)
	assert.Equal(t, num, uint64(100))
	_ = clpKeeper.GetParams(ctx)
	_ = clpKeeper.Logger(ctx)
	pool.ExternalAsset.Symbol = ""
	err := clpKeeper.SetPool(ctx, &pool)
	assert.Error(t, err, "Unable to set pool")
	boolean := clpKeeper.ValidatePool(pool)
	assert.False(t, boolean)
	getpools, _, err := clpKeeper.GetPoolsPaginated(ctx, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(getpools), 0, "No pool added")
	lp := test.GenerateRandomLP(1)[0]
	lp.Asset.Symbol = ""
	clpKeeper.SetLiquidityProvider(ctx, lp)
	getlp, err := clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.Error(t, err)
	assert.NotEqual(t, getlp, lp)
	assert.NotNil(t, test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA7"))
}

func TestKeeper_BankKeeper(t *testing.T) {
	user1 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	user2 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA7")
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	initialBalance := sdk.NewUint(10000)
	sendingBalance := sdk.NewUint(1000)
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(initialBalance))
	sendingCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sendingBalance))
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, user1, sdk.NewCoins(nativeCoin))
	assert.NoError(t, err)
	assert.True(t, clpKeeper.HasBalance(ctx, user1, nativeCoin))
	assert.NoError(t, clpKeeper.SendCoins(ctx, user1, user2, sdk.NewCoins(sendingCoin)))
	assert.True(t, clpKeeper.HasBalance(ctx, user2, sendingCoin))
}

func TestKeeper_GetModuleAccount(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	moduleAccount := clpKeeper.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	assert.Equal(t, moduleAccount.GetName(), types.ModuleName)
	assert.Equal(t, moduleAccount.GetPermissions(), []string{authtypes.Burner, authtypes.Minter})
}

func TestKeeper_Codec(t *testing.T) {
	_, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	got := clpKeeper.Codec()
	require.NotNil(t, got)
}

func TestKeeper_GetBankKeeper(t *testing.T) {
	_, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	got := clpKeeper.GetBankKeeper()
	require.NotNil(t, got)
}

// nolint
func TestKeeper_GetAssetDecimals(t *testing.T) {
	testcases :=
		[]struct {
			name        string
			denom       string
			asset       types.Asset
			decimals    int64
			createToken bool
			errString   error
			expected    uint8
		}{
			{
				name:        "big decimals number throws error",
				asset:       types.Asset{Symbol: "xxx"},
				createToken: true,
				denom:       "xxx",
				decimals:    256,
				errString:   errors.New("Could not perform type cast"),
			},
			{
				name:        "negative decimals number throws error",
				asset:       types.Asset{Symbol: "xxx"},
				createToken: true,
				denom:       "xxx",
				decimals:    -200,
				errString:   errors.New("Could not perform type cast"),
			},
			{
				name:        "unknown symbol",
				createToken: false,
				asset:       types.Asset{Symbol: "xxx"},
				errString:   errors.New("registry entry not found: key not found"),
			},
			{
				name:        "success",
				asset:       types.Asset{Symbol: "xxx"},
				createToken: true,
				denom:       "xxx",
				decimals:    73,
				expected:    73,
			},
		}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClp(false)
			clpKeeper := app.ClpKeeper

			if tc.createToken {
				app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
					Denom:       tc.denom,
					Decimals:    tc.decimals,
					Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
				})
			}
			decimals, err := clpKeeper.GetAssetDecimals(ctx, tc.asset)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			require.Equal(t, tc.expected, decimals)
		})
	}
}
