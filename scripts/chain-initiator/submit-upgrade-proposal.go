package main

import (
	"log"
	"os/exec"
	"strings"
)

func submitUpgradeProposal(cmdPath, name, newVersion, upgradeHeight, homePath, keyringBackend, chainId, node, broadcastMode string) {
	planName := newVersion
	// Remove the "v" prefix if present
	if strings.HasPrefix(planName, "v") {
		planName = strings.TrimPrefix(planName, "v")
	}

	// Command and arguments
	args := []string{
		"tx",
		"gov",
		// "submit-legacy-proposal", // not available in v0.45.x
		"submit-proposal",
		"software-upgrade",
		planName,
		"--title", newVersion,
		"--description", newVersion,
		"--upgrade-height", upgradeHeight,
		// "--no-validate", // not available in v0.45.x
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "5000000000000000000000rowan",
		"--gas", "1000000",
		"--deposit", "50000000000000000000000rowan",
		"--home", homePath,
		"--yes",
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"Submitted upgrade proposal: %s, upgrade block height: %s", newVersion, upgradeHeight)
}
