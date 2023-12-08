package main

import (
	"log"
	"os/exec"
)

func genTx(cmdPath, name, amount, chainId, homePath, keyringBackend string) {
	// Command and arguments
	args := []string{"gentx", name, amount + "rowan", "--chain-id", chainId, "--home", homePath, "--keyring-backend", keyringBackend}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"gen tx with name %s, amount: %s, chain id %s, home path %s and keyring backend %s successfully", name, amount, chainId, homePath, keyringBackend)
}
