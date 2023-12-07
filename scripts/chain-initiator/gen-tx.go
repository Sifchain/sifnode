package main

import (
	"log"
	"os/exec"
)

func genTx(cmdPath, name, balance, chainId, homePath, keyringBackend string) {
	// Command and arguments
	args := []string{"gentx", name, balance + "rowan", "--chain-id", chainId, "--home", homePath, "--keyring-backend", keyringBackend}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf("gen tx with name %s, balance: %s, chain id %s, home path %s and keyring backend %s successfully", name, balance, chainId, homePath, keyringBackend)
}
