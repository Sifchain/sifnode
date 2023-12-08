package main

import (
	"log"
	"os/exec"
)

func retrieveSnapshot(snapshotUrl, homePath string) {
	// Construct the command string
	cmdString := "curl -o - -L " + snapshotUrl + " | lz4 -c -d - | tar -x -C " + homePath

	// Execute the command using /bin/sh
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	if err := cmd.Run(); err != nil {
		log.Fatalf(Red+"Command execution failed: %v", err)
	}

	// If execution reaches here, the command was successful
	log.Printf(Yellow+"Snapshot retrieved and extracted to path: %s", homePath)
}
