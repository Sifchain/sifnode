package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func writeGenesisFile(filePath string, genesis Genesis) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	// encoder.SetIndent("", "  ") // disable for now

	if err := encoder.Encode(genesis); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}
