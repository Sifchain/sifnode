package main

import (
	"log"
	"strconv"
)

func queryAndCalcUpgradeBlockHeight(cmdPath, node string) string {
	// query block height
	blockHeight := queryBlockHeight(cmdPath, node)

	// Convert blockHeight from string to int
	blockHeightInt, err := strconv.Atoi(blockHeight)
	if err != nil {
		log.Fatalf("Failed to convert blockHeight to integer: %v", err)
	}

	// set upgrade block height
	upgradeBlockHeight := blockHeightInt + 100

	// return upgrade block height as a string
	return strconv.Itoa(upgradeBlockHeight)
}
