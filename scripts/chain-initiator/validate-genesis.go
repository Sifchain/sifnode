package main

import (
	"log"
	"os/exec"
)

func validateGenesis(cmdPath, homePath string) {
	// Command and arguments
	args := []string{"validate-genesis", "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"validate genesis with home path %s successfully", homePath)
}
