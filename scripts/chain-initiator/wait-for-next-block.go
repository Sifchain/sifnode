package main

import (
	"log"
	"strconv"
	"time"
)

func waitForNextBlock(cmdPath, node string) {
	var currentBlockHeight, newBlockHeight int
	var err error

	// First, get the current block height
	for {
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			currentBlockHeight, err = strconv.Atoi(blockHeightStr)
			if err == nil && currentBlockHeight > 0 {
				break
			}
		}
		time.Sleep(5 * time.Second) // Wait 5 second before retrying
	}

	log.Printf(Yellow+"Current Block Height: %d", currentBlockHeight)

	// Now, wait for the block height to increase
	for {
		var blockHeightStr string
		blockHeightStr, err = queryBlockHeight(cmdPath, node)
		if err == nil {
			newBlockHeight, err = strconv.Atoi(blockHeightStr)
			if err == nil && newBlockHeight > currentBlockHeight {
				break
			}
		}
		time.Sleep(5 * time.Second) // Wait a second before retrying
	}

	log.Printf(Yellow+"New Block Height: %d", newBlockHeight)
}
