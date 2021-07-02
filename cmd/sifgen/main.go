package main

import (
	"github.com/Sifchain/sifnode/cmd/sifgen/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	_ = rootCmd.Execute()
}
