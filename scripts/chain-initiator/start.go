package main

import (
	"log"
	"os"
	"os/exec"
)

func start(cmdPath, homePath string) {
	// Command and arguments
	args := []string{"start", "--home", homePath}

	// Set up the command
	cmd := exec.Command(cmdPath, args...)

	// Attach command's stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and stream the output
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
