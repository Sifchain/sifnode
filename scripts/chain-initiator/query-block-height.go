package main

import (
	"encoding/json"
	"log"
	"os/exec"
)

func queryBlockHeight(cmdPath, node string) string {
	// Command and arguments
	args := []string{"status", "--node", node}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	// Unmarshal the JSON output
	var statusOutput StatusOutput
	if err := json.Unmarshal(output, &statusOutput); err != nil {
		log.Fatalf("Failed to unmarshal JSON output: %v", err)
	}

	return statusOutput.SyncInfo.LatestBlockHeight
}
