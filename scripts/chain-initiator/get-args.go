package main

import (
	"log"
)

func getArgs(args []string) (snapshotUrl, oldBinaryUrl, newBinaryUrl string) {
	snapshotUrl = args[0] // https://snapshots.polkachu.com/snapshots/sifchain/sifchain_15048938.tar.lz4
	if snapshotUrl == "" {
		log.Fatalf(Red + "snapshot url is required")
	}

	oldBinaryUrl = args[1] // https://github.com/Sifchain/sifnode/releases/download/v1.2.0-beta/sifnoded-v1.2.0-beta-darwin-arm64
	if oldBinaryUrl == "" {
		log.Fatalf(Red + "old binary url is required")
	}

	newBinaryUrl = args[2] // https://github.com/Sifchain/sifnode/releases/download/v1.3.0-beta/sifnoded-v1.3.0-beta-darwin-arm64
	if newBinaryUrl == "" {
		log.Fatalf(Red + "new binary url is required")
	}

	return
}
