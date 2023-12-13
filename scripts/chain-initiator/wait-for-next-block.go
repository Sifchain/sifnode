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
		log.Println(Yellow + "Waiting for current block height...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
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
		log.Println(Yellow + "Waiting for next block height...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	log.Printf(Yellow+"New Block Height: %d", newBlockHeight)
}
