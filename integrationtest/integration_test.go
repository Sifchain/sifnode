package integrationtest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	marginkeeper "github.com/Sifchain/sifnode/x/margin/keeper"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/proto/tendermint/types"
)

type TestCase struct {
	Name  string
	Setup struct {
		Accounts         []banktypes.Balance
		Margin           *margintypes.GenesisState
		RewardsParams    clptypes.RewardParams
		ProtectionParams clptypes.LiquidityProtectionParams
		ShiftingParams   clptypes.PmtpParams
		ProviderParams   clptypes.ProviderDistributionParams
	}
	Messages []sdk.Msg
}

func TC1(t *testing.T) TestCase {
	externalAsset := "cusdc"
	address := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

	externalAssetBalance, ok := sdk.NewIntFromString("1000000000000000")
	require.True(t, ok)
	nativeAssetBalance, ok := sdk.NewIntFromString("1000000000000000000000000000")
	require.True(t, ok)
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin(externalAsset, externalAssetBalance),
				sdk.NewCoin("rowan", nativeAssetBalance),
			},
		},
	}

	tc := TestCase{
		Setup: struct {
			Accounts         []banktypes.Balance
			Margin           *margintypes.GenesisState
			RewardsParams    clptypes.RewardParams
			ProtectionParams clptypes.LiquidityProtectionParams
			ShiftingParams   clptypes.PmtpParams
			ProviderParams   clptypes.ProviderDistributionParams
		}{
			Accounts:       balances,
			ShiftingParams: *clptypes.GetDefaultPmtpParams(),
		},
		Messages: []sdk.Msg{
			&clptypes.MsgCreatePool{
				Signer:              address,
				ExternalAsset:       &clptypes.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
				ExternalAssetAmount: sdk.NewUintFromString("1000000000"),             // 1000cusdc
			},
			&clptypes.MsgAddLiquidity{
				Signer:              address,
				ExternalAsset:       &clptypes.Asset{Symbol: externalAsset},
				NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
				ExternalAssetAmount: sdk.NewUintFromString("1000000000"),
			},
			&margintypes.MsgOpen{
				Signer:           address,
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUintFromString("10000000000000000000"), // 10rowan
				BorrowAsset:      externalAsset,
				Position:         margintypes.Position_LONG,
				Leverage:         sdk.NewDec(2),
			},
			&clptypes.MsgSwap{
				Signer:             address,
				SentAsset:          &clptypes.Asset{Symbol: externalAsset},
				ReceivedAsset:      &clptypes.Asset{Symbol: clptypes.NativeSymbol},
				SentAmount:         sdk.NewUintFromString("10000"),
				MinReceivingAmount: sdk.NewUint(0),
			},
			&clptypes.MsgRemoveLiquidity{
				Signer:        address,
				ExternalAsset: &clptypes.Asset{Symbol: externalAsset},
				WBasisPoints:  sdk.NewInt(5000),
				Asymmetry:     sdk.NewInt(0),
			},
			&margintypes.MsgClose{
				Signer: address,
				Id:     1,
			},
		},
	}

	return tc
}

func TestIntegration(t *testing.T) {
	overwriteFlag := flag.Bool("overwrite", false, "Overwrite test output")
	flag.Parse()

	tt := []TestCase{
		TC1(t),
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				// Initialise token registry
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: "cusdc", BaseDenom: "cusdc", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				// Initialise wallet
				bankGs := banktypes.DefaultGenesisState()
				bankGs.Balances = append(bankGs.Balances, tc.Setup.Accounts...)
				bz, _ = app.AppCodec().MarshalJSON(bankGs)
				genesisState["bank"] = bz

				// Set enabled margin pools
				marginGs := margintypes.DefaultGenesis()
				marginGs.Params.Pools = append(marginGs.Params.Pools, "cusdc")
				bz, _ = app.AppCodec().MarshalJSON(marginGs)
				genesisState["margin"] = bz

				return genesisState
			})

			app.ClpKeeper.SetRewardParams(ctx, &tc.Setup.RewardsParams)
			app.ClpKeeper.SetLiquidityProtectionParams(ctx, &tc.Setup.ProtectionParams)
			app.ClpKeeper.SetPmtpParams(ctx, &tc.Setup.ShiftingParams)
			app.ClpKeeper.SetProviderDistributionParams(ctx, &tc.Setup.ProviderParams)

			clpSrv := clpkeeper.NewMsgServerImpl(app.ClpKeeper)
			marginSrv := marginkeeper.NewMsgServerImpl(app.MarginKeeper)

			for i, msg := range tc.Messages {
				ctx = ctx.WithBlockHeight(int64(i))
				app.BeginBlocker(ctx, abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})
				switch msg := msg.(type) {
				case *clptypes.MsgCreatePool:
					_, err := clpSrv.CreatePool(sdk.WrapSDKContext(ctx), msg)
					require.NoError(t, err)
				case *clptypes.MsgAddLiquidity:
					_, err := clpSrv.AddLiquidity(sdk.WrapSDKContext(ctx), msg)
					require.NoError(t, err)
				case *clptypes.MsgRemoveLiquidity:
					_, err := clpSrv.RemoveLiquidity(sdk.WrapSDKContext(ctx), msg)
					require.NoError(t, err)
				case *clptypes.MsgSwap:
					_, err := clpSrv.Swap(sdk.WrapSDKContext(ctx), msg)
					require.NoError(t, err)
				case *margintypes.MsgOpen:
					_, err := marginSrv.Open(sdk.WrapSDKContext(ctx), msg)
					require.NoError(t, err)
				case *margintypes.MsgClose:
					_, err := marginSrv.Close(sdk.WrapSDKContext(ctx), msg)
					require.NoError(t, err)
				}
				endBlock(t, app, ctx, ctx.BlockHeight())
			}

			// Check balances
			results := getResults(t, app, ctx, tc.Setup.Accounts)
			if *overwriteFlag {
				writeResults(t, results)
			} else {
				expected, err := getExpected()
				require.NoError(t, err)
				require.EqualValues(t, expected, &results)
			}
		})
	}

}

func endBlock(t *testing.T, app *sifapp.SifchainApp, ctx sdk.Context, height int64) {
	app.EndBlocker(ctx, abci.RequestEndBlock{Height: height})
	//app.Commit()

	// Check invariants
	res, stop := app.ClpKeeper.BalanceModuleAccountCheck()(ctx)
	require.False(t, stop, res)
}

type TestResults struct {
	Accounts map[string]sdk.Coins `json:"accounts"`
	Pools    map[string]clptypes.Pool
	LPs      map[string]clptypes.LiquidityProvider
}

func getResults(t *testing.T, app *sifapp.SifchainApp, ctx sdk.Context, accounts []banktypes.Balance) TestResults {
	pools := app.ClpKeeper.GetPools(ctx)

	lps, err := app.ClpKeeper.GetAllLiquidityProviders(ctx)
	require.NoError(t, err)

	results := TestResults{
		Accounts: make(map[string]sdk.Coins, len(accounts)),
		Pools:    make(map[string]clptypes.Pool, len(pools)),
		LPs:      make(map[string]clptypes.LiquidityProvider, len(lps)),
	}

	for _, account := range accounts {
		// Lookup account balances
		addr, err := sdk.AccAddressFromBech32(account.Address)
		require.NoError(t, err)
		balances := app.BankKeeper.GetAllBalances(ctx, addr)
		results.Accounts[account.Address] = balances
	}

	for _, pool := range pools {
		results.Pools[pool.ExternalAsset.Symbol] = *pool
	}

	for _, lp := range lps {
		results.LPs[string(clptypes.GetLiquidityProviderKey(lp.Asset.Symbol, lp.LiquidityProviderAddress))] = *lp
	}

	return results
}

func writeResults(t *testing.T, results TestResults) {
	bz, err := json.Marshal(results)
	fmt.Printf("%s", bz)
	require.NoError(t, err)
	err = os.WriteFile("output/results.json", bz, 0644)
	require.NoError(t, err)
}

func getExpected() (*TestResults, error) {
	bz, err := os.ReadFile("output/results.json")
	if err != nil {
		return nil, err
	}
	var results TestResults
	err = json.Unmarshal(bz, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
