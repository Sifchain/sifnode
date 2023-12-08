package main

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	flagHome = "home"
	flagCmd  = "cmd"
)

func getFlags(cmd *cobra.Command) (homePath, cmdPath string) {
	homePath, _ = cmd.Flags().GetString(flagHome)
	if homePath == "" {
		log.Fatalf("home path is required")
	}

	cmdPath, _ = cmd.Flags().GetString(flagCmd)
	if cmdPath == "" {
		log.Fatalf("cmd path is required")
	}

	return
}
