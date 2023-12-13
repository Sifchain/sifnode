package main

import (
	"log"
	"os/exec"
)

func addGenesisAccount(cmdPath, address, balance, homePath string) {
	// Command and arguments
	args := []string{"add-genesis-account", address, balance + "rowan", "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err) // nolint: goconst
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"add genesis account with address %s, balance: %s and home path %s successfully", address, balance, homePath)
}
