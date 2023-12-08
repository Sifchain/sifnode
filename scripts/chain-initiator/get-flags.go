package main

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	flagHome          = "home"
	flagCmd           = "cmd"
	flagSkipSnapshot  = "skip-snapshot"
	flagSkipChainInit = "skip-chain-init"
	flagSkipNodeStart = "skip-node-start"
)

func getFlags(cmd *cobra.Command) (homePath, cmdPath string, skipSnapshot, skipChainInit, skipNodeStart bool) {
	homePath, _ = cmd.Flags().GetString(flagHome)
	if homePath == "" {
		log.Fatalf(Red + "home path is required")
	}

	cmdPath, _ = cmd.Flags().GetString(flagCmd)
	if cmdPath == "" {
		log.Fatalf(Red + "cmd path is required")
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

	return
}
