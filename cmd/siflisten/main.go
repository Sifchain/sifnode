package main

import (
	"github.com/Sifchain/sifnode/app"
	chainevents "github.com/Sifchain/sifnode/tools/siflisten/chain/events"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	cmd := &cobra.Command{
		Use: "siflisten",
	}

	syncCmd := syncCmd()

	cmd.AddCommand(syncCmd)

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func syncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync events from chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeURI, err := cmd.Flags().GetString("node")
			if err != nil {
				return err
			}

			rpcClient, err := rpchttp.New(nodeURI, "/websocket")
			if err != nil {
				return err
			}

			encConfig := app.MakeTestEncodingConfig()

			clientCtx := sdkclient.Context{}.
				//WithNodeURI(uri).
				WithClient(rpcClient).
				//WithFrom(args.Sender.String()).
				//WithFromAddress(args.Sender).
				//WithTxConfig(config.TxConfig).
				WithInterfaceRegistry(encConfig.InterfaceRegistry).
				//WithSkipConfirmation(true).
				//WithFromName(args.SenderName).
				//WithBroadcastMode("block").
				WithOutputFormat("json")

			chainevents.Sync(clientCtx, nil)

			return nil
		}}

	cmd.PersistentFlags().String("node", "tcp://127.0.0.1:26657", "Tendermint RPC node to watch events on")

	return cmd
}
