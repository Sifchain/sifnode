package main

import (
	"encoding/json"
	"log"
	"os/exec"
)

func addKey(cmdPath, name, homePath, keyringBackend string) string {
	// Command and arguments
	args := []string{"keys", "add", name, "--home", homePath, "--keyring-backend", keyringBackend, "--output", "json"}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	// Unmarshal the JSON output
	var keyOutput KeyOutput
	if err := json.Unmarshal(output, &keyOutput); err != nil {
		log.Fatalf("Failed to unmarshal JSON output: %v", err)
	}

	// Log the address
	log.Printf("add key with name %s, home path: %s, keyring backend %s and address %s successfully", name, homePath, keyringBackend, keyOutput.Address)

	return keyOutput.Address
}
