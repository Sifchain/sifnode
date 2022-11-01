package oracle_test

import (
	"bytes"
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

//nolint:golint
func TestInitGenesis(t *testing.T) {
	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)

	tt, _ := testGenesisData(t)

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			keeper.RemoveOracleWhiteList(ctx, networkDescriptor)

			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)

			if len(tc.genesis.AdminAddress) <= 0 {
				require.Nil(t, keeper.GetAdminAccount(ctx))
			} else {
				require.Equal(t, tc.genesis.AdminAddress, keeper.GetAdminAccount(ctx).String())
			}

			wl := keeper.GetOracleWhiteList(ctx, networkDescriptor).ValidatorPower

			whiteList := tc.genesis.NetworkConfigData

			found := false
			if len(wl) == 0 {
				found = true
			}
			expectedWhitelist := make([]*types.ValidatorPower, 0)

			for _, value := range whiteList {
				if value.NetworkDescriptor == networkDescriptor.NetworkDescriptor {
					found = true
					expectedWhitelist = value.ValidatorWhitelist.ValidatorPower
				}
			}
			assert.Equal(t, found, true)

			for _, value := range wl {
				found := false
				for _, expected := range expectedWhitelist {
					if bytes.Compare(value.ValidatorAddress, expected.ValidatorAddress) == 0 {
						found = true
						assert.Equal(t, value.VotingPower, expected.VotingPower)
					}
				}
				assert.Equal(t, found, true)
			}

			prophecies := keeper.GetProphecies(ctx)
			require.Equal(t, len(tc.genesis.Prophecies), len(prophecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, *p, prophecies[i])
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	tt, _ := testGenesisData(t)

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			networkDescriptor := types.NetworkIdentity{NetworkDescriptor: types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM}

			keeper.RemoveOracleWhiteList(ctx, networkDescriptor)

			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)
			genesis := oracle.ExportGenesis(ctx, keeper)
			require.Equal(t, tc.genesis.AdminAddress, genesis.AdminAddress)

			wl := genesis.GetNetworkConfigData()
			// require.Equal(t, len(tc.genesis.NetworkConfigData), len(wl))

			found := false
			expectedWhitelist := make([]*types.ValidatorPower, 0)

			if len(wl) == 0 {
				found = true
			}

			for _, value := range wl {
				if value.NetworkDescriptor == networkDescriptor.NetworkDescriptor {
					found = true
					if value.ValidatorWhitelist != nil {
						expectedWhitelist = value.ValidatorWhitelist.ValidatorPower
					}

				}
			}
			assert.Equal(t, found, true)

			found = false
			if len(tc.genesis.NetworkConfigData) == 0 {
				found = true
			}
			genesisWhitelist := make([]*types.ValidatorPower, 0)
			for _, value := range tc.genesis.NetworkConfigData {
				if value.NetworkDescriptor == networkDescriptor.NetworkDescriptor {
					found = true
					genesisWhitelist = value.ValidatorWhitelist.ValidatorPower
				}
			}
			assert.Equal(t, found, true)

			for _, value := range genesisWhitelist {
				found := false
				for _, expected := range expectedWhitelist {
					if bytes.Compare(value.ValidatorAddress, expected.ValidatorAddress) == 0 {
						found = true
						assert.Equal(t, value.VotingPower, expected.VotingPower)
					}
				}
				assert.Equal(t, found, true)
			}

			prophecies := genesis.Prophecies
			require.Equal(t, len(tc.genesis.Prophecies), len(prophecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, p, prophecies[i])
			}
		})
	}
}

func TestGenesisMarshalling(t *testing.T) {
	tt, prophecies := testGenesisData(t)

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, keeper, encCfg, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
			networkDescriptor := types.NetworkIdentity{NetworkDescriptor: types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM}
			keeper.RemoveOracleWhiteList(ctx, networkDescriptor)

			_ = oracle.InitGenesis(ctx, keeper, tc.genesis)
			genesis := oracle.ExportGenesis(ctx, keeper)

			genesisData := encCfg.Marshaler.MustMarshalJSON(genesis)

			var genesisState types.GenesisState
			encCfg.Marshaler.MustUnmarshalJSON(genesisData, &genesisState)

			ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 1, []int64{1}, "")

			_ = oracle.InitGenesis(ctx, keeper, genesisState)

			require.Equal(t, tc.genesis.AdminAddress, genesis.AdminAddress)

			wl := genesis.NetworkConfigData

			found := false
			expectedWhitelist := make([]*types.ValidatorPower, 0)

			if len(wl) == 0 {
				found = true
			}

			for _, value := range wl {
				if value.NetworkDescriptor == networkDescriptor.NetworkDescriptor {
					found = true
					if value.ValidatorWhitelist != nil {
						expectedWhitelist = value.ValidatorWhitelist.ValidatorPower
					}

				}
			}
			assert.Equal(t, found, true)

			found = false
			if len(tc.genesis.NetworkConfigData) == 0 {
				found = true
			}
			genesisWhitelist := make([]*types.ValidatorPower, 0)
			for _, value := range tc.genesis.NetworkConfigData {
				if value.NetworkDescriptor == networkDescriptor.NetworkDescriptor {
					found = true
					genesisWhitelist = value.ValidatorWhitelist.ValidatorPower
				}
			}
			assert.Equal(t, found, true)

			for _, value := range genesisWhitelist {
				found := false
				for _, expected := range expectedWhitelist {
					if bytes.Compare(value.ValidatorAddress, expected.ValidatorAddress) == 0 {
						found = true
						assert.Equal(t, value.VotingPower, expected.VotingPower)
					}
				}
				assert.Equal(t, found, true)
			}

			dbProphecies := genesis.Prophecies
			require.Equal(t, len(tc.genesis.Prophecies), len(dbProphecies))
			for i, p := range tc.genesis.Prophecies {
				require.Equal(t, p, dbProphecies[i])
				require.Equal(t, prophecies[i], *p)
			}
		})
	}
}

type testCase struct {
	name    string
	genesis types.GenesisState
}

func testGenesisData(t *testing.T) ([]testCase, []types.Prophecy) {
	addrs, valAddrs := test.CreateTestAddrs(2)
	// power := uint32(100)

	whitelist := make([]string, len(valAddrs))
	for i, addr := range valAddrs {
		whitelist[i] = addr.String()
	}

	prophecy := types.Prophecy{
		Id:              []byte("asd"),
		Status:          types.StatusText_STATUS_TEXT_PENDING,
		ClaimValidators: []string{valAddrs[0].String()},
	}

	return []testCase{
		{
			name:    "Default genesis",
			genesis: *types.DefaultGenesisState(),
		},
		{
			name: "Prophecy",
			genesis: types.GenesisState{
				NetworkConfigData: []*types.NetworkConfigData{},
				AdminAddress:      addrs[0].String(),
				Prophecies: []*types.Prophecy{
					&prophecy,
				},
			},
		},
	}, []types.Prophecy{prophecy}
}

func TestGenesisWithCrossChainFee(t *testing.T) {
	ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
	networkIdentity := types.NewNetworkIdentity(types.NetworkDescriptor(1))
	one := sdk.NewInt(1)
	keeper.SetCrossChainFee(ctx, networkIdentity, "rowan", one, one, one, one)

	exportedState := oracle.ExportGenesis(ctx, keeper)
	newCtx, _, _, _, newKeeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")

	oracle.InitGenesis(newCtx, newKeeper, *exportedState)

	assert.Equal(t, keeper.GetAllCrossChainFeeConfig(ctx), newKeeper.GetAllCrossChainFeeConfig(newCtx))
}

func TestGenesisWithConsensusNeeded(t *testing.T) {
	ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
	networkIdentity := types.NewNetworkIdentity(types.NetworkDescriptor(1))
	keeper.SetConsensusNeeded(ctx, networkIdentity, types.ConsensusNeeded{ConsensusNeeded: 66})

	exportedState := oracle.ExportGenesis(ctx, keeper)
	newCtx, _, _, _, newKeeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")

	oracle.InitGenesis(newCtx, newKeeper, *exportedState)

	assert.Equal(t, keeper.GetAllCrossChainFeeConfig(ctx), newKeeper.GetAllCrossChainFeeConfig(newCtx))
}

func TestGenesisWithWitnessLockBurnSequence(t *testing.T) {
	ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")
	_, valAddrs := test.CreateTestAddrs(2)

	networkDescriptor := types.NetworkDescriptor(1)
	keeper.SetWitnessLockBurnSequence(ctx, networkDescriptor, valAddrs[0], 66)

	exportedState := oracle.ExportGenesis(ctx, keeper)
	newCtx, _, _, _, newKeeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")

	oracle.InitGenesis(newCtx, newKeeper, *exportedState)

	assert.Equal(t, keeper.GetAllWitnessLockBurnSequence(ctx), newKeeper.GetAllWitnessLockBurnSequence(newCtx))
}

func TestGenesisWithProphecyInfo(t *testing.T) {
	ctx, _, _, _, keeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")

	prophecyID := []byte{1, 2, 3, 4, 5, 6}
	networkDescriptor := types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	cosmosSender := "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	cosmosSenderSequence := uint64(1)
	ethereumReceiver := "0x0000000000000000000000000000000000000003"
	tokenDenomHash := "rowan"
	tokenContractAddress := "0x0000000000000000000000000000000000000004"
	tokenAmount := sdk.NewInt(1)
	crosschainFee := sdk.NewInt(1)
	bridgeToken := true
	globalSequence := uint64(1)
	tokenDecimal := uint8(1)
	tokenName := "name"
	tokenSymbol := "symbol"

	err := keeper.SetProphecyInfo(ctx, prophecyID,
		networkDescriptor,
		cosmosSender,
		cosmosSenderSequence,
		ethereumReceiver,
		tokenDenomHash,
		tokenContractAddress,
		tokenAmount,
		crosschainFee,
		bridgeToken,
		globalSequence,
		tokenDecimal,
		tokenName,
		tokenSymbol)

	assert.NoError(t, err)

	exportedState := oracle.ExportGenesis(ctx, keeper)
	newCtx, _, _, _, newKeeper, _, _, _ := test.CreateTestKeepers(t, 1, []int64{1}, "")

	oracle.InitGenesis(newCtx, newKeeper, *exportedState)

	exportedProphecyID, existed := newKeeper.GetProphecyIDByNetworkDescriptorGlobalNonce(newCtx, networkDescriptor, globalSequence)
	assert.Equal(t, existed, true)
	assert.Equal(t, prophecyID, exportedProphecyID)

	assert.Equal(t, keeper.GetAllProphecyInfo(ctx), newKeeper.GetAllProphecyInfo(newCtx))
}
