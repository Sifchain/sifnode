package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func filterBalances(balances []banktypes.Balance, filterAddresses []string) ([]banktypes.Balance, sdk.Coins) {
	filterMap := make(map[string]struct{})
	for _, addr := range filterAddresses {
		filterMap[addr] = struct{}{}
	}

	newBalances := []banktypes.Balance{}
	var coinsToRemove sdk.Coins
	for _, balance := range balances {
		if _, exists := filterMap[balance.Address]; exists {
			coinsToRemove = coinsToRemove.Add(balance.Coins...)
			continue
		}
		newBalances = append(newBalances, balance)
	}
	return newBalances, coinsToRemove
}
