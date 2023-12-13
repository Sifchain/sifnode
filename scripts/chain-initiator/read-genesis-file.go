package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func readGenesisFile(filePath string) (Genesis, error) {
	var genesis Genesis
	file, err := os.Open(filePath)
	if err != nil {
		return genesis, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(bufio.NewReader(file)).Decode(&genesis); err != nil {
		return genesis, fmt.Errorf("error decoding JSON: %w", err)
	}

	return genesis, nil
}
