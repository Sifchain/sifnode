package v42

import (
	v039clp "github.com/Sifchain/sifnode/x/clp/legacy/v39"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
)

func Migrate(genesis v039clp.GenesisState) clptypes.GenesisState {
	whitelist := make([]string, len(genesis.AddressWhitelist))
	for i, addr := range genesis.AddressWhitelist {
		whitelist[i] = addr.String()
	}

	poolList := make([]*clptypes.Pool, len(genesis.PoolList))
	for i, pool := range genesis.PoolList {
		poolList[i] = &clptypes.Pool{
			ExternalAsset:        &clptypes.Asset{Symbol: pool.ExternalAsset.Symbol},
			NativeAssetBalance:   pool.NativeAssetBalance,
			ExternalAssetBalance: pool.ExternalAssetBalance,
			PoolUnits:            pool.PoolUnits,
		}
	}

	liquidityProviders := make([]*clptypes.LiquidityProvider, len(genesis.LiquidityProviderList))
	for i, lp := range genesis.LiquidityProviderList {
		liquidityProviders[i] = &clptypes.LiquidityProvider{
			Asset:                    &clptypes.Asset{Symbol: lp.Asset.Symbol},
			LiquidityProviderUnits:   lp.LiquidityProviderUnits,
			LiquidityProviderAddress: lp.LiquidityProviderAddress.String(),
		}
	}

	return clptypes.GenesisState{
		Params:             clptypes.Params{MinCreatePoolThreshold: uint64(genesis.Params.MinCreatePoolThreshold)},
		AddressWhitelist:   whitelist,
		PoolList:           poolList,
		LiquidityProviders: liquidityProviders,
	}
}
