// nolint: nakedret
package main

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	flagHome                    = "home"
	flagSkipSnapshot            = "skip-snapshot"
	flagSkipChainInit           = "skip-chain-init"
	flagSkipNodeStart           = "skip-node-start"
	flagSkipProposal            = "skip-proposal"
	flagSkipBinary              = "skip-binary"
	flagMoniker                 = "moniker"
	flagChainId                 = "chain-id"
	flagKeyringBackend          = "keyring-backend"
	flagValidatorKeyName        = "validator-key-name"
	flagValidatorBalance        = "validator-balance"
	flagValidatorSelfDelegation = "validator-self-delegation"
	flagGenesisFilePath         = "genesis-file-path"
	flagNode                    = "node"
	flagBroadcastMode           = "broadcast-mode"
)

func getFlags(cmd *cobra.Command) (homePath string, skipSnapshot, skipChainInit, skipNodeStart, skipProposal, skipBinary bool, moniker, chainId, keyringBackend, validatorKeyName, validatorBalance, validatorSelfDelegation, genesisFilePath, node, broadcastMode string) {
	homePath, _ = cmd.Flags().GetString(flagHome)
	if homePath == "" {
		log.Fatalf(Red + "home path is required")
	}

	skipSnapshot, _ = cmd.Flags().GetBool(flagSkipSnapshot)
	if skipSnapshot {
		log.Printf(Yellow + "skipping snapshot retrieval")
	}

	skipChainInit, _ = cmd.Flags().GetBool(flagSkipChainInit)
	if skipChainInit {
		log.Printf(Yellow + "skipping chain init")
	}

	skipNodeStart, _ = cmd.Flags().GetBool(flagSkipNodeStart)
	if skipNodeStart {
		log.Printf(Yellow + "skipping node start")
	}

	skipProposal, _ = cmd.Flags().GetBool(flagSkipProposal)
	if skipProposal {
		log.Printf(Yellow + "skipping proposal")
	}

	skipBinary, _ = cmd.Flags().GetBool(flagSkipBinary)
	if skipBinary {
		log.Printf(Yellow + "skipping binary download")
	}

	moniker, _ = cmd.Flags().GetString(flagMoniker)
	if moniker == "" {
		log.Fatalf(Red + "moniker is required")
	}

	chainId, _ = cmd.Flags().GetString(flagChainId)
	if chainId == "" {
		log.Fatalf(Red + "chain id is required")
	}

	keyringBackend, _ = cmd.Flags().GetString(flagKeyringBackend)
	if keyringBackend == "" {
		log.Fatalf(Red + "keyring backend is required")
	}

	validatorKeyName, _ = cmd.Flags().GetString(flagValidatorKeyName)
	if validatorKeyName == "" {
		log.Fatalf(Red + "validator key name is required")
	}

	validatorBalance, _ = cmd.Flags().GetString(flagValidatorBalance)
	if validatorBalance == "" {
		log.Fatalf(Red + "validator balance is required")
	}

	validatorSelfDelegation, _ = cmd.Flags().GetString(flagValidatorSelfDelegation)
	if validatorSelfDelegation == "" {
		log.Fatalf(Red + "validator self delegation is required")
	}

	genesisFilePath, _ = cmd.Flags().GetString(flagGenesisFilePath)
	if genesisFilePath == "" {
		log.Fatalf(Red + "genesis file path is required")
	}

	node, _ = cmd.Flags().GetString(flagNode)
	if node == "" {
		log.Fatalf(Red + "node is required")
	}

	broadcastMode, _ = cmd.Flags().GetString(flagBroadcastMode)
	if broadcastMode == "" {
		log.Fatalf(Red + "broadcast mode is required")
	}

	return
}
