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
	ethtest "github.com/Sifchain/sifnode/x/ethbridge/test"
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
		Name: "tc1",
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

func TC2(t *testing.T) TestCase {
	sifapp.SetConfig(false)
	externalAsset := "cusdc"
	address := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	addresses, _ := ethtest.CreateTestAddrs(1)

	externalAssetBalance, ok := sdk.NewIntFromString("1000000000000000000") // 1,000,000,000,000.000000
	require.True(t, ok)
	nativeAssetBalance, ok := sdk.NewIntFromString("1000000000000000000000000000000") // 1000,000,000,000.000000000000000000
	require.True(t, ok)
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("atom", externalAssetBalance),
				sdk.NewCoin(externalAsset, externalAssetBalance),
				sdk.NewCoin("rowan", nativeAssetBalance),
			},
		},
		{
			Address: addresses[0].String(),
			Coins: sdk.Coins{
				sdk.NewCoin("atom", externalAssetBalance),
				sdk.NewCoin(externalAsset, externalAssetBalance),
				sdk.NewCoin("rowan", nativeAssetBalance),
			},
		},
	}
	allocation := sdk.NewUintFromString("1000000000000000000000000")
	defaultMultiplier := sdk.NewDec(1)

	tc := TestCase{
		Name: "tc2",
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
			RewardsParams: clptypes.RewardParams{
				LiquidityRemovalLockPeriod:   0,
				LiquidityRemovalCancelPeriod: 0,
				RewardPeriods: []*clptypes.RewardPeriod{
					&clptypes.RewardPeriod{
						RewardPeriodId:                "1",
						RewardPeriodStartBlock:        1,
						RewardPeriodEndBlock:          1000,
						RewardPeriodAllocation:        &allocation,
						RewardPeriodPoolMultipliers:   []*clptypes.PoolMultiplier{},
						RewardPeriodDefaultMultiplier: &defaultMultiplier,
						RewardPeriodDistribute:        false,
						RewardPeriodMod:               1,
					},
				},
				RewardPeriodStartTime: "",
			},
			ProviderParams: clptypes.ProviderDistributionParams{
				DistributionPeriods: []*clptypes.ProviderDistributionPeriod{
					&clptypes.ProviderDistributionPeriod{
						DistributionPeriodBlockRate:  sdk.NewDecWithPrec(7, 6),
						DistributionPeriodStartBlock: 1,
						DistributionPeriodEndBlock:   1000,
						DistributionPeriodMod:        1,
					},
				},
			},
		},
		Messages: []sdk.Msg{
			&clptypes.MsgCreatePool{
				Signer:              address,
				ExternalAsset:       &clptypes.Asset{Symbol: "atom"},
				NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000000000"), // 1000,000,000rowan
				ExternalAssetAmount: sdk.NewUintFromString("1000000000000000"),             // 1000,000,000atom
			},
			&margintypes.MsgOpen{
				Signer:           address,
				CollateralAsset:  "atom",
				CollateralAmount: sdk.NewUintFromString("500000000"), // 500atom
				BorrowAsset:      "rowan",
				Position:         margintypes.Position_LONG,
				Leverage:         sdk.NewDec(10),
			},
			&margintypes.MsgOpen{
				Signer:           addresses[0].String(),
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUintFromString("500000000000000000000000"), // 500,000rowan
				BorrowAsset:      "atom",
				Position:         margintypes.Position_LONG,
				Leverage:         sdk.NewDec(5),
			}, /*
				&clptypes.MsgAddLiquidity{
					Signer:              address,
					ExternalAsset:       &clptypes.Asset{Symbol: externalAsset},
					NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
					ExternalAssetAmount: sdk.NewUintFromString("1000000000"),
				},*/
		},
	}

	return tc
}

func TC3(t *testing.T) TestCase {
	sifapp.SetConfig(false)
	externalAsset := "cusdc"
	address := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	addresses, _ := ethtest.CreateTestAddrs(2)

	externalAssetBalance, ok := sdk.NewIntFromString("1000000000000000000") // 1,000,000,000,000.000000
	require.True(t, ok)
	nativeAssetBalance, ok := sdk.NewIntFromString("1000000000000000000000000000000") // 1000,000,000,000.000000000000000000
	require.True(t, ok)
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("atom", externalAssetBalance),
				sdk.NewCoin(externalAsset, externalAssetBalance),
				sdk.NewCoin("rowan", nativeAssetBalance),
			},
		},
		{
			Address: addresses[0].String(),
			Coins: sdk.Coins{
				sdk.NewCoin("atom", externalAssetBalance),
				sdk.NewCoin(externalAsset, externalAssetBalance),
				sdk.NewCoin("rowan", nativeAssetBalance),
			},
		},
		{
			Address: addresses[1].String(),
			Coins: sdk.Coins{
				sdk.NewCoin("atom", externalAssetBalance),
				sdk.NewCoin(externalAsset, externalAssetBalance),
				sdk.NewCoin("rowan", nativeAssetBalance),
			},
		},
	}
	allocation := sdk.NewUintFromString("2000000000000000000000000")
	defaultMultiplier := sdk.NewDec(1)

	tc := TestCase{
		Name: "tc3",
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
			RewardsParams: clptypes.RewardParams{
				LiquidityRemovalLockPeriod:   0,
				LiquidityRemovalCancelPeriod: 0,
				RewardPeriods: []*clptypes.RewardPeriod{
					&clptypes.RewardPeriod{
						RewardPeriodId:                "1",
						RewardPeriodStartBlock:        1,
						RewardPeriodEndBlock:          1000,
						RewardPeriodAllocation:        &allocation,
						RewardPeriodPoolMultipliers:   []*clptypes.PoolMultiplier{},
						RewardPeriodDefaultMultiplier: &defaultMultiplier,
						RewardPeriodDistribute:        false,
						RewardPeriodMod:               1,
					},
				},
				RewardPeriodStartTime: "",
			},
			ProviderParams: clptypes.ProviderDistributionParams{
				DistributionPeriods: []*clptypes.ProviderDistributionPeriod{
					&clptypes.ProviderDistributionPeriod{
						DistributionPeriodBlockRate:  sdk.NewDecWithPrec(7, 6),
						DistributionPeriodStartBlock: 1,
						DistributionPeriodEndBlock:   1000,
						DistributionPeriodMod:        1,
					},
				},
			},
		},
		Messages: []sdk.Msg{
			&clptypes.MsgCreatePool{
				Signer:              address,
				ExternalAsset:       &clptypes.Asset{Symbol: "atom"},
				NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000000000"), // 1000,000,000rowan
				ExternalAssetAmount: sdk.NewUintFromString("1000000000000000"),             // 1000,000,000atom
			},
			&clptypes.MsgCreatePool{
				Signer:              address,
				ExternalAsset:       &clptypes.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000000000"), // 1000,000,000rowan
				ExternalAssetAmount: sdk.NewUintFromString("1000000000000000"),             // 1000,000,000cusdc
			},
			&margintypes.MsgOpen{
				Signer:           address,
				CollateralAsset:  "cusdc",
				CollateralAmount: sdk.NewUintFromString("1000000000"), // 1000atom
				BorrowAsset:      "rowan",
				Position:         margintypes.Position_LONG,
				Leverage:         sdk.NewDec(10),
			},
			&margintypes.MsgOpen{
				Signer:           addresses[0].String(),
				CollateralAsset:  "cusdc",
				CollateralAmount: sdk.NewUintFromString("5000000000"), // 5000
				BorrowAsset:      "rowan",
				Position:         margintypes.Position_LONG,
				Leverage:         sdk.NewDec(5),
			},
			&margintypes.MsgOpen{
				Signer:           addresses[1].String(),
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUintFromString("500000000000000000000000"), // 500,000
				BorrowAsset:      "atom",
				Position:         margintypes.Position_LONG,
				Leverage:         sdk.NewDec(3),
			},
			&clptypes.MsgSwap{
				Signer:             address,
				SentAsset:          &clptypes.Asset{Symbol: "atom"},
				ReceivedAsset:      &clptypes.Asset{Symbol: clptypes.NativeSymbol},
				SentAmount:         sdk.NewUintFromString("5000000000"), // 5000
				MinReceivingAmount: sdk.NewUint(0),
			}, /*
				&clptypes.MsgAddLiquidity{
					Signer:              address,
					ExternalAsset:       &clptypes.Asset{Symbol: externalAsset},
					NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
					ExternalAssetAmount: sdk.NewUintFromString("1000000000"),
				},*/
		},
	}

	return tc
}

func TestIntegration(t *testing.T) {
	overwriteFlag := flag.Bool("overwrite", false, "Overwrite test output")
	flag.Parse()

	tt := []TestCase{
		TC1(t), TC2(t), TC3(t),
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				// Initialise token registry
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: "atom", BaseDenom: "atom", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
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
				marginGs.Params.Pools = append(marginGs.Params.Pools, []string{"cusdc", "atom"}...)
				bz, _ = app.AppCodec().MarshalJSON(marginGs)
				genesisState["margin"] = bz

				return genesisState
			})

			app.ClpKeeper.SetRewardParams(ctx, &tc.Setup.RewardsParams)                 //nolint
			app.ClpKeeper.SetLiquidityProtectionParams(ctx, &tc.Setup.ProtectionParams) //nolint
			app.ClpKeeper.SetPmtpParams(ctx, &tc.Setup.ShiftingParams)                  //nolint
			app.ClpKeeper.SetProviderDistributionParams(ctx, &tc.Setup.ProviderParams)  //nolint

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
			results := getResults(t, app, ctx, tc)
			if *overwriteFlag {
				writeResults(t, tc, results)
			} else {
				expected, err := getExpected(tc)
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

func getResults(t *testing.T, app *sifapp.SifchainApp, ctx sdk.Context, tc TestCase) TestResults {
	pools := app.ClpKeeper.GetPools(ctx)

	lps, err := app.ClpKeeper.GetAllLiquidityProviders(ctx)
	require.NoError(t, err)

	results := TestResults{
		Accounts: make(map[string]sdk.Coins, len(tc.Setup.Accounts)),
		Pools:    make(map[string]clptypes.Pool, len(pools)),
		LPs:      make(map[string]clptypes.LiquidityProvider, len(lps)),
	}

	for _, account := range tc.Setup.Accounts {
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

func writeResults(t *testing.T, tc TestCase, results TestResults) {
	bz, err := json.MarshalIndent(results, "", "\t")
	fmt.Printf("%s", bz)
	require.NoError(t, err)

	filename := "output/" + tc.Name + ".json"

	err = os.WriteFile(filename, bz, 0600)
	require.NoError(t, err)
}

func getExpected(tc TestCase) (*TestResults, error) {
	bz, err := os.ReadFile("output/" + tc.Name + ".json")
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
