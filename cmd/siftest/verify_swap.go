package main

import (
	"fmt"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetVerifySwap() *cobra.Command {
	- Target environment
	- Block height
	- Sent symbol
	- Sent amount
	- Received symbol
	- Slippage
	- # of swaps

	cmd := &cobra.Command{
		Use:   "swap --env --height --sent_symbol --sent_amount --received_symbol --slippage --number_of_swap",
		Short: "Verify a removal",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("verifying removal...\n")
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			unitsRemoved := sdk.NewUintFromString(viper.GetString("units"))

			//err = VerifySwap()
			if err != nil {
				panic(err)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	//cmd.Flags().Uint64("height", 0, "height of transaction")
	cmd.Flags().String("from", "", "address of transactor")
	cmd.Flags().String("units", "0", "number of units removed")
	cmd.Flags().String("external-asset", "", "external asset of pool")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("units")
	cmd.MarkFlagRequired("external-asset")
	cmd.MarkFlagRequired("height")
	return cmd
}


//- Environment
//- Block height
//- Swap rate ROWAN:TKN
//- Swap rate TKN:TKN
//- Swap rate TKN:ROWAN
//- Sent amount
//- Intermediate amount
//- ROWAN value if TKN:TKN swap
//- empty if ROWAN:TKN or TKN: ROWAN swap
//- Received amount
//- Before wallet balance (sent symbol)
//- Before wallet balanced (received symbol)
//- After wallet balance (sent symbol)
//- After wallet balance (received symbol)
//- Pool depth at time of transaction
//- Pool custody amount
//- Pool liability amount
func VerifySwap(clientCtx client.Context,height uint64,
	sent_symbol string ,
    sent_amount sdk.Uint,
	received_symbol string,
    slippage sdk.Dec,
    number_of_swap uint64 ) error {
    clpQueryClient := clptypes.NewQueryClient(clientCtx.WithHeight(int64(height)))
	clpQueryClient.GetPool()

	return nil

}