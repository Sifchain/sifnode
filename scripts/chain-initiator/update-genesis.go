package main

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func updateGenesis(validatorBalance, homePath, genesisFilePath string) {
	genesis, err := readGenesisFile(genesisFilePath)
	if err != nil {
		log.Fatalf(Red+"Error reading genesis file: %v", err)
	}

	genesisInitFilePath := homePath + "/config/genesis.json"
	genesisInit, err := readGenesisFile(genesisInitFilePath)
	if err != nil {
		log.Fatalf(Red+"Error reading initial genesis file: %v", err)
	}

	filterAccountAddresses := []string{
		"sif1harggtyrlukcfrtmpgjzptsnaedcdh38qqknp2", // multisig account with missing pubkeys
	}
	filterBalanceAddresses := []string{
		"sif1harggtyrlukcfrtmpgjzptsnaedcdh38qqknp2",
		authtypes.NewModuleAddress("distribution").String(),
		authtypes.NewModuleAddress("bonded_tokens_pool").String(),
		authtypes.NewModuleAddress("not_bonded_tokens_pool").String(),
	}

	var coinsToRemove sdk.Coins

	genesis.AppState.Auth.Accounts = filterAccounts(genesis.AppState.Auth.Accounts, filterAccountAddresses)
	genesis.AppState.Bank.Balances, coinsToRemove = filterBalances(genesis.AppState.Bank.Balances, filterBalanceAddresses)

	newValidatorBalance, ok := sdk.NewIntFromString(validatorBalance)
	if !ok {
		panic(Red + "invalid number")
	}
	newValidatorBalanceCoin := sdk.NewCoin("rowan", newValidatorBalance)

	// update supply
	genesis.AppState.Bank.Supply = genesis.AppState.Bank.Supply.Sub(coinsToRemove).Add(newValidatorBalanceCoin)

	// Add new validator account and balance
	genesis.AppState.Auth.Accounts = append(genesis.AppState.Auth.Accounts, genesisInit.AppState.Auth.Accounts[0])
	genesis.AppState.Bank.Balances = append(genesis.AppState.Bank.Balances, genesisInit.AppState.Bank.Balances[0])

	// reset staking data
	stakingParams := genesis.AppState.Staking.Params
	genesis.AppState.Staking = genesisInit.AppState.Staking
	genesis.AppState.Staking.Params = stakingParams

	// reset slashing data
	genesis.AppState.Slashing = genesisInit.AppState.Slashing

	// reset distribution data
	genesis.AppState.Distribution = genesisInit.AppState.Distribution

	// set genutil from genesisInit
	genesis.AppState.Genutil = genesisInit.AppState.Genutil

	// update voting period
	genesis.AppState.Gov.VotingParams.VotingPeriod = "10s"
	genesis.AppState.Gov.DepositParams.MaxDepositPeriod = "10s"

	outputFilePath := homePath + "/config/genesis.json"
	if err := writeGenesisFile(outputFilePath, genesis); err != nil {
		log.Fatalf(Red+"Error writing genesis file: %v", err)
	}
}
