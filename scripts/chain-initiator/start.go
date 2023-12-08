package main

import (
	"log"
	"os"
	"os/exec"
)

func start(cmdPath, homePath string) *exec.Cmd {
	// Command and arguments
	args := []string{"start", "--home", homePath}

	// Set up the command
	cmd := exec.Command(cmdPath, args...)

	// Attach command's stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and stream the output in a goroutine to avoid blocking
	go func() {
		if err := cmd.Run(); err != nil {
			log.Fatalf("Command execution failed: %v", err)
		}
	}()

	return cmd
}
