package main

import (
	"log"

	ebrelayer "github.com/Sifchain/sifnode/cmd/ebrelayer"
	sifgen "github.com/Sifchain/sifnode/cmd/sifgen"
	sifnoded "github.com/Sifchain/sifnode/cmd/sifnoded/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	sifnodedCmd, _ := sifnoded.NewRootCmd()
	sifgenCmd := sifgen.NewRootCmd()
	ebrelayerCmd := ebrelayer.NewRootCmd()

	err := doc.GenMarkdownTree(sifnodedCmd, "docs/cmd/sifnoded/")
	if err != nil {
		log.Fatal(err)
	}

	err = doc.GenMarkdownTree(sifgenCmd, "docs/cmd/sifgen/")
	if err != nil {
		log.Fatal(err)
	}

	err = doc.GenMarkdownTree(ebrelayerCmd, "docs/cmd/ebrelayer/")
	if err != nil {
		log.Fatal(err)
	}
}
