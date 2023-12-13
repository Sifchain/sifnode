package main

import (
	"log"
	"os/exec"
)

func initChain(cmdPath, moniker, chainId, homePath string) {
	// Command and arguments
	args := []string{"init", moniker, "--chain-id", chainId, "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"init chain with moniker %s, chain id %s and home path: %s successfully", moniker, chainId, homePath)
}
