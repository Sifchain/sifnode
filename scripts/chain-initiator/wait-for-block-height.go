package main

import (
	"log"
	"strconv"
	"time"
)

func waitForBlockHeight(cmdPath, node, height string) {
	targetBlockHeight, err := strconv.Atoi(height)
	if err != nil {
		log.Fatalf(Red+"Error converting target block height to integer: %v", err)
	}

	// Now, wait for the block height
	for {
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			newBlockHeight, err := strconv.Atoi(blockHeightStr)
			if err == nil && newBlockHeight >= targetBlockHeight {
				break
			}
		}
		log.Println(Yellow+"Waiting for block height", height, "...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(Yellow+"Block height %d reached", targetBlockHeight)
}
