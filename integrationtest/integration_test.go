package integrationtest

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	externalAsset := "cusdc"
	address := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

	tt := []struct {
		name          string
		externalAsset string
		address       string
		createPoolMsg clptypes.MsgCreatePool
	}{
		{
			externalAsset: externalAsset,
			address:       address,
			createPoolMsg: clptypes.MsgCreatePool{
				Signer:              address,
				ExternalAsset:       &clptypes.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
				ExternalAssetAmount: sdk.NewUintFromString("1000000000"),             // 1000cusdc
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				// Initialise token registry
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: externalAsset, BaseDenom: externalAsset, Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				// Initialise wallet
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
				bankGs := banktypes.DefaultGenesisState()
				bankGs.Balances = append(bankGs.Balances, balances...)
				bz, _ = app.AppCodec().MarshalJSON(bankGs)
				genesisState["bank"] = bz

				return genesisState
			})

			clpSrv := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			_, err := clpSrv.CreatePool(sdk.WrapSDKContext(ctx), &tc.createPoolMsg)
			require.NoError(t, err)

			// TODO: Check invariants
		})
	}

}
