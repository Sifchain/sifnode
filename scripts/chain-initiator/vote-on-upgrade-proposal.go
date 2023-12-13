package main

import (
	"log"
	"os/exec"
)

func voteOnUpgradeProposal(cmdPath, name, proposalId, homePath, keyringBackend, chainId, node, broadcastMode string) {
	// Command and arguments
	args := []string{
		"tx", "gov", "vote", proposalId, "yes",
		"--from", name,
		"--keyring-backend", keyringBackend,
		"--chain-id", chainId,
		"--node", node,
		"--broadcast-mode", broadcastMode,
		"--fees", "100000000000000000rowan",
		"--gas", "1000000",
		"--home", homePath,
		"--yes",
	}

	// Execute the command
	if err := exec.Command(cmdPath, args...).Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"Voted on upgrade proposal: %s", proposalId)
}
