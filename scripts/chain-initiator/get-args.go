package main

import (
	"log"
)

func getArgs(args []string) (snapshotUrl, newVersion string) {
	snapshotUrl = args[0] // https://snapshots.polkachu.com/snapshots/sifchain/sifchain_15048938.tar.lz4
	if snapshotUrl == "" {
		log.Fatalf(Red + "snapshot url is required")
	}

	newVersion = args[1] // v0.1.0
	if newVersion == "" {
		log.Fatalf(Red + "new version is required")
	}

	return
}
