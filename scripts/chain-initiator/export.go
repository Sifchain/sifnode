package main

import (
	"io/ioutil"
	"log"
	"os/exec"
)

func export(cmdPath, homePath, genesisFilePath string) {
	// Command and arguments
	args := []string{"export", "--home", homePath}

	// Execute the command and capture the output
	output, err := exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// Write the output to the specified file
	err = ioutil.WriteFile(genesisFilePath, output, 0644) // nolint: gosec
	if err != nil {
		log.Fatalf(Red+"Failed to write output to file: %v", err)
	}

	log.Printf(Yellow+"Output successfully written to %s", genesisFilePath)
}
