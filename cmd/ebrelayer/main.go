package main

import (
	"os"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/cmd"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), sifapp.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
