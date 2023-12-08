package main

import (
	"log"
	"os"

	app "github.com/Sifchain/sifnode/app"
	"github.com/spf13/cobra"
)

const (
	moniker                 = "node"
	chainId                 = "sifchain-1"
	keyringBackend          = "test"
	validatorKeyName        = "validator"
	validatorBalance        = "4000000000000000000000000000"
	validatorSelfDelegation = "1000000000000000000000000000"
	genesisFilePath         = "/tmp/genesis.json"
	node                    = "tcp://localhost:26657"
	broadcastMode           = "block"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "initiator [snapshot_url] [new_version] [flags]",
		Short: "Chain Initiator is a tool for running a chain from a snapshot.",
		Long:  `A tool for running a chain from a snapshot.`,
		Args:  cobra.ExactArgs(2), // Expect exactly 1 argument
		Run: func(cmd *cobra.Command, args []string) {
			snapshotUrl, newVersion := getArgs(args)
			homePath, cmdPath, skipSnapshot, skipChainInit, skipNodeStart := getFlags(cmd)

			// set address prefix
			app.SetConfig(false)

			if !skipSnapshot {
				// remove home path
				removeHome(homePath)

				// init chain
				initChain(cmdPath, moniker, chainId, homePath)

				// retrieve the snapshot
				retrieveSnapshot(snapshotUrl, homePath)

				// export genesis file
				export(cmdPath, homePath, genesisFilePath)
			}

			if !skipChainInit {
				// remove home path
				removeHome(homePath)

				// init chain
				initChain(cmdPath, moniker, chainId, homePath)

				// add validator key
				validatorAddress := addKey(cmdPath, validatorKeyName, homePath, keyringBackend)

				// add genesis account
				addGenesisAccount(cmdPath, validatorAddress, validatorBalance, homePath)

				// generate genesis tx
				genTx(cmdPath, validatorKeyName, validatorSelfDelegation, chainId, homePath, keyringBackend)

				// collect genesis txs
				collectGentxs(cmdPath, homePath)

				// validate genesis
				validateGenesis(cmdPath, homePath)

				// update genesis
				updateGenesis(validatorBalance, homePath)
			}

			if !skipNodeStart {
				// start chain
				startCmd := start(cmdPath, homePath)

				// wait for node to start
				waitForNodeToStart(node)

				// wait for next block
				waitForNextBlock(cmdPath, node)

				// query and calculate upgrade block height
				upgradeBlockHeight := queryAndCalcUpgradeBlockHeight(cmdPath, node)

				// submit upgrade proposal
				submitUpgradeProposal(cmdPath, validatorKeyName, newVersion, upgradeBlockHeight, homePath, keyringBackend, chainId, node, broadcastMode)

				// listen for signals
				listenForSignals(startCmd)
			}
		},
	}

	// get HOME environment variable
	homeEnv, _ := os.LookupEnv("HOME")

	rootCmd.PersistentFlags().String(flagCmd, homeEnv+"/go/bin/sifnoded", "path to sifnoded")
	rootCmd.PersistentFlags().String(flagHome, homeEnv+"/.sifnoded", "home directory")
	rootCmd.PersistentFlags().Bool(flagSkipSnapshot, false, "skip snapshot retrieval")
	rootCmd.PersistentFlags().Bool(flagSkipChainInit, false, "skip chain init")
	rootCmd.PersistentFlags().Bool(flagSkipNodeStart, false, "skip node start")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(Red+"Error executing command: %v", err)
	}
}
