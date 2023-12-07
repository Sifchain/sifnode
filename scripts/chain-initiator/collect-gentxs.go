package main

import (
	"log"
	"os/exec"
)

func collectGentxs(cmdPath, homePath string) {
	// Command and arguments
	args := []string{"collect-gentxs", "--home", homePath}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf("collect gen txs with home path %s successfully", homePath)
}
