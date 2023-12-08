package main

import (
	"encoding/json"
	"os/exec"
)

func queryBlockHeight(cmdPath, node string) (string, error) {
	// Command and arguments
	args := []string{"status", "--node", node}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		return "-1", err
	}

	// Unmarshal the JSON output
	var statusOutput StatusOutput
	if err := json.Unmarshal(output, &statusOutput); err != nil {
		return "-1", err
	}

	return statusOutput.SyncInfo.LatestBlockHeight, nil
}
