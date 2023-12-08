package main

import (
	"log"

	app "github.com/Sifchain/sifnode/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

const (
	moniker          = "node"
	chainId          = "sifchain-1"
	keyringBackend   = "test"
	validatorKeyName = "validator"
	validatorBalance = "4000000000000000000000000000"
	genesisFilePath  = "/tmp/genesis.json"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "initiator [cmd_path] [home_path] [snapshot_url]]",
		Short: "Chain Initiator is a tool for modifying genesis files",
		Long:  `A tool for performing various operations on genesis files of a blockchain setup.`,
		Args:  cobra.ExactArgs(3), // Expect exactly two arguments
		Run: func(cmd *cobra.Command, args []string) {
			cmdPath := args[0]     // sifnoded
			homePath := args[1]    // /tmp/node
			snapshotUrl := args[2] // https://snapshots.polkachu.com/snapshots/sifchain/sifchain_15048938.tar.lz4

			// set address prefix
			app.SetConfig(false)

			// remove home path
			removeHome(homePath)

			// init chain
			initChain(cmdPath, moniker, chainId, homePath)

			// retrieve the snapshot
			retrieveSnapshot(snapshotUrl, homePath)

			// export genesis file
			export(cmdPath, homePath, genesisFilePath)

			// remove home path
			removeHome(homePath)

			// init chain
			initChain(cmdPath, moniker, chainId, homePath)

			// add validator key
			validatorAddress := addKey(cmdPath, validatorKeyName, homePath, keyringBackend)

			// add genesis account
			addGenesisAccount(cmdPath, validatorAddress, validatorBalance, homePath)

			// generate genesis tx
			genTx(cmdPath, validatorKeyName, validatorBalance, chainId, homePath, keyringBackend)

			// collect genesis txs
			collectGentxs(cmdPath, homePath)

			// validate genesis
			validateGenesis(cmdPath, homePath)

			genesis, err := readGenesisFile(genesisFilePath)
			if err != nil {
				log.Fatalf("Error reading genesis file: %v", err)
			}

			genesisInitFilePath := homePath + "/config/genesis.json"
			genesisInit, err := readGenesisFile(genesisInitFilePath)
			if err != nil {
				log.Fatalf("Error reading initial genesis file: %v", err)
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

			newValidatorBalance, ok := sdk.NewIntFromString("4000000000000000000000000000")
			if !ok {
				panic("invalid number")
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

			outputFilePath := homePath + "/config/genesis.json"
			if err := writeGenesisFile(outputFilePath, genesis); err != nil {
				log.Fatalf("Error writing genesis file: %v", err)
			}

			// start chain
			start(cmdPath, homePath)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
