package main

import (
	"log"
	"os"
	"time"

	app "github.com/Sifchain/sifnode/app"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "initiator [snapshot_url] [old_binary_url] [new_binary_url] [flags]",
		Short: "Chain Initiator is a tool for running a chain from a snapshot.",
		Long:  `A tool for running a chain from a snapshot.`,
		Args:  cobra.ExactArgs(3), // Expect exactly 1 argument
		Run: func(cmd *cobra.Command, args []string) {
			snapshotUrl, oldBinaryUrl, newBinaryUrl := getArgs(args)
			homePath, skipSnapshot, skipChainInit, skipNodeStart, skipProposal, skipBinary, moniker, chainId, keyringBackend, validatorKeyName, validatorBalance, validatorSelfDelegation, genesisFilePath, node, broadcastMode := getFlags(cmd)

			// set address prefix
			app.SetConfig(false)

			// download and run old binary
			oldBinaryPath, oldVersion, err := downloadAndRunVersion(oldBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(Red+"Error downloading and running old binary: %v", err)
			}

			// print old binary path and version
			log.Printf(Green+"Old binary path: %v and version: %v", oldBinaryPath, oldVersion)

			// download and run new binary
			newBinaryPath, newVersion, err := downloadAndRunVersion(newBinaryUrl, skipBinary)
			if err != nil {
				log.Fatalf(Red+"Error downloading and running new binary: %v", err)
			}

			// print new binary path and version
			log.Printf(Green+"New binary path: %v and version: %v", newBinaryPath, newVersion)

			if !skipSnapshot {
				// remove home path
				removeHome(homePath)

				// init chain
				initChain(oldBinaryPath, moniker, chainId, homePath)

				// retrieve the snapshot
				retrieveSnapshot(snapshotUrl, homePath)

				// export genesis file
				export(oldBinaryPath, homePath, genesisFilePath)
			}

			if !skipChainInit {
				// remove home path
				removeHome(homePath)

				// init chain
				initChain(oldBinaryPath, moniker, chainId, homePath)

				// add validator key
				validatorAddress := addKey(oldBinaryPath, validatorKeyName, homePath, keyringBackend)

				// add genesis account
				addGenesisAccount(oldBinaryPath, validatorAddress, validatorBalance, homePath)

				// generate genesis tx
				genTx(oldBinaryPath, validatorKeyName, validatorSelfDelegation, chainId, homePath, keyringBackend)

				// collect genesis txs
				collectGentxs(oldBinaryPath, homePath)

				// validate genesis
				validateGenesis(oldBinaryPath, homePath)

				// update genesis
				updateGenesis(validatorBalance, homePath, genesisFilePath)
			}

			if !skipNodeStart {
				// start chain
				oldBinaryCmd := start(oldBinaryPath, homePath)

				// wait for node to start
				waitForNodeToStart(node)

				// wait for next block
				waitForNextBlock(oldBinaryPath, node)

				if skipProposal {
					// listen for signals
					listenForSignals(oldBinaryCmd)
					return
				}

				// query and calculate upgrade block height
				upgradeBlockHeight := queryAndCalcUpgradeBlockHeight(oldBinaryPath, node)

				// query next proposal id
				proposalId, err := queryNextProposalId(oldBinaryPath, node)
				if err != nil {
					log.Fatalf(Red+"Error querying next proposal id: %v", err)
				}

				// submit upgrade proposal
				submitUpgradeProposal(oldBinaryPath, validatorKeyName, newVersion, upgradeBlockHeight, homePath, keyringBackend, chainId, node, broadcastMode)

				// vote on upgrade proposal
				voteOnUpgradeProposal(oldBinaryPath, validatorKeyName, proposalId, homePath, keyringBackend, chainId, node, broadcastMode)

				// wait for upgrade block height
				waitForBlockHeight(oldBinaryPath, node, upgradeBlockHeight)

				// stop old binary
				stop(oldBinaryCmd)

				// wait 5 seconds
				time.Sleep(5 * time.Second)

				// start new binary
				newBinaryCmd := start(newBinaryPath, homePath)

				// wait for node to start
				waitForNodeToStart(node)

				// wait for next block
				waitForNextBlock(newBinaryPath, node)

				// listen for signals
				listenForSignals(newBinaryCmd)
			}
		},
	}

	// get HOME environment variable
	homeEnv, _ := os.LookupEnv("HOME")

	rootCmd.PersistentFlags().String(flagHome, homeEnv+"/.sifnoded", "home directory")
	rootCmd.PersistentFlags().Bool(flagSkipSnapshot, false, "skip snapshot retrieval")
	rootCmd.PersistentFlags().Bool(flagSkipChainInit, false, "skip chain init")
	rootCmd.PersistentFlags().Bool(flagSkipNodeStart, false, "skip node start")
	rootCmd.PersistentFlags().Bool(flagSkipProposal, false, "skip proposal")
	rootCmd.PersistentFlags().Bool(flagSkipBinary, false, "skip binary download")
	rootCmd.PersistentFlags().String(flagMoniker, "node", "moniker")
	rootCmd.PersistentFlags().String(flagChainId, "sifchain-1", "chain id")
	rootCmd.PersistentFlags().String(flagKeyringBackend, "test", "keyring backend")
	rootCmd.PersistentFlags().String(flagValidatorKeyName, "validator", "validator key name")
	rootCmd.PersistentFlags().String(flagValidatorBalance, "4000000000000000000000000000", "validator balance")
	rootCmd.PersistentFlags().String(flagValidatorSelfDelegation, "1000000000000000000000000000", "validator self delegation")
	rootCmd.PersistentFlags().String(flagGenesisFilePath, "/tmp/genesis.json", "genesis file path")
	rootCmd.PersistentFlags().String(flagNode, "tcp://localhost:26657", "node")
	rootCmd.PersistentFlags().String(flagBroadcastMode, "block", "broadcast mode")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(Red+"Error executing command: %v", err)
	}
}
