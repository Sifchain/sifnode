package clp

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	types2 "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	abci "github.com/tendermint/tendermint/abci/types"
	"io/ioutil"
	"math"
	"path/filepath"
)

func GeneratePerfData(ctx sdk.Context, k keeper.Keeper) error {
	var inputs types.Pools
	file, err := filepath.Abs("pools.json")
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(input, &inputs)
	if err != nil {
		return err
	}
	k.GetTokenRegistryKeeper().SetRegistry(ctx, types2.Registry{})
	multipliers := []*types.PoolMultiplier{}
	cethMultiplier := sdk.MustNewDecFromStr("1.5")
	rp_1_allocation := sdk.NewUintFromString("1550459183129248235861408")
	for _, pool := range inputs {
		err := k.SetPool(ctx, &pool)
		if err != nil {
			return err
		}
		r := types2.RegistryEntry{
			Decimals:                 18,
			Denom:                    pool.ExternalAsset.Symbol,
			BaseDenom:                "",
			Path:                     "",
			IbcChannelId:             "",
			IbcCounterpartyChannelId: "",
			DisplayName:              "",
			DisplaySymbol:            "",
			Network:                  "",
			Address:                  "",
			ExternalSymbol:           "",
			TransferLimit:            "",
			Permissions:              []types2.Permission{types2.Permission_CLP},
			UnitDenom:                "",
			IbcCounterpartyDenom:     "",
			IbcCounterpartyChainId:   "",
		}
		multipliers = append(multipliers, &types.PoolMultiplier{
			PoolMultiplierAsset: pool.ExternalAsset.Symbol,
			Multiplier:          &cethMultiplier,
		})
		k.GetTokenRegistryKeeper().SetToken(ctx, &r)
	}
	rp := types.RewardParams{
		LiquidityRemovalLockPeriod:   1,
		LiquidityRemovalCancelPeriod: 70,
		RewardPeriods: []*types.RewardPeriod{
			{
				RewardPeriodId:                "RP_1",
				RewardPeriodStartBlock:        1,
				RewardPeriodEndBlock:          1100,
				RewardPeriodAllocation:        &rp_1_allocation,
				RewardPeriodPoolMultipliers:   multipliers,
				RewardPeriodDefaultMultiplier: &cethMultiplier,
			},
		},
	}
	k.SetRewardParams(ctx, &rp)
	return nil
}
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	err := GeneratePerfData(ctx, k)
	if err != nil {
		fmt.Println(err)
	}
	k.SetParams(ctx, data.Params)
	//k.SetRewardParams(ctx, types.GetDefaultRewardParams())
	// Initiate Pmtp
	k.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.ZeroDec(),
		PmtpCurrentRunningRate: sdk.ZeroDec(),
		PmtpInterPolicyRate:    sdk.ZeroDec(),
	})
	k.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: 0,
		BlockCounter: 0,
	})
	k.SetPmtpParams(ctx, types.GetDefaultPmtpParams())

	k.SetPmtpInterPolicyRate(ctx, sdk.NewDec(0))
	if data.AddressWhitelist == nil || len(data.AddressWhitelist) == 0 {
		panic("AddressWhiteList must be set.")
	}
	wl := make([]sdk.AccAddress, len(data.AddressWhitelist))
	if data.AddressWhitelist != nil {
		for i, entry := range data.AddressWhitelist {
			wlAddress, err := sdk.AccAddressFromBech32(entry)
			if err != nil {
				panic(err)
			}
			wl[i] = wlAddress
		}
		k.SetClpWhiteList(ctx, wl)
	}
	k.SetClpWhiteList(ctx, wl)
	for _, pool := range data.PoolList {
		err := k.SetPool(ctx, pool)
		if err != nil {
			panic(fmt.Sprintf("Pool could not be set : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviders {
		k.SetLiquidityProvider(ctx, lp)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	var poolList []*types.Pool
	poolList, _, _ = keeper.GetPoolsPaginated(ctx, &query.PageRequest{
		Limit: uint64(math.MaxUint64),
	})
	liquidityProviders, _, _ := keeper.GetAllLiquidityProvidersPaginated(ctx, &query.PageRequest{
		Limit: uint64(math.MaxUint64),
	})
	whiteList := keeper.GetClpWhiteList(ctx)
	wl := make([]string, len(whiteList))
	for i, entry := range whiteList {
		wl[i] = entry.String()
	}
	return types.GenesisState{
		Params:             params,
		AddressWhitelist:   wl,
		PoolList:           poolList,
		LiquidityProviders: liquidityProviders,
	}
}

// ValidateGenesis validates the clp genesis parameters
func ValidateGenesis(data types.GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("clp: params are invalid : %s \n %s", err.Error(), data.Params.String()))
	}
	for _, pool := range data.PoolList {
		if !pool.Validate() {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("clp: pool is invalid : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviders {
		if !lp.Validate() {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("clp: liquidityProvider is invalid : %s", lp.String()))
		}
	}
	return nil
}
