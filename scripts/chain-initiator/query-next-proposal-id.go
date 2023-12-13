package main

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
)

func queryNextProposalId(cmdPath, node string) (string, error) {
	// Command and arguments
	args := []string{"q", "gov", "proposals", "--node", node, "--limit", "1", "--reverse", "--output", "json"}

	// Execute the command
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		return "-1", err
	}

	// Unmarshal the JSON output
	var proposalsOutput ProposalsOutput
	if err := json.Unmarshal(output, &proposalsOutput); err != nil {
		return "-1", err
	}

	// check if there are any proposals
	if len(proposalsOutput.Proposals) == 0 {
		return "1", errors.New("no proposals found")
	}

	// increment proposal id
	proposalId := proposalsOutput.Proposals[0].ProposalId
	proposalIdInt, err := strconv.Atoi(proposalId)
	if err != nil {
		return "-1", err
	}
	proposalIdInt++
	// convert back to string
	proposalId = strconv.Itoa(proposalIdInt)

	return proposalId, nil
}
