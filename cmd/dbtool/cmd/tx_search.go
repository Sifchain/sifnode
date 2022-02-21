/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Sifchain/sifnode/cmd/dbtool/utils"

	"github.com/spf13/cobra"
)

var (
	outfile = fmt.Sprintf("%s/dbtool.data", homeDir())
	full    = false
	pages   = 1
	perPage = 10
)

func NewSearchCmd() *cobra.Command {
	txSearchCmd := &cobra.Command{
		Use:   "tx-search",
		Short: "Search for transactions by event criteria",
		Long: `Search for transactions by event criteria
		
ex: 
dbtool tx-search "message.action='/ibc.core.client.v1.MsgUpdateClient'"
dbtool tx-search "fungible_token_packet.denom='ujuno'" 
dbtool tx-search "fungible_token_packet.sucess='001'"
dbtool tx-search "fungible_token_packet.denom='transfer/channel-19/ungm'"
dbtool tx-search "fungible_token_packet.denom='transfer/channel-18/uluna'"
dbtool tx-search "update_client.client_id='07-tendermint-41' AND fungible_token_packet.success='false'"
dbtool tx-search "update_client.client_id='07-tendermint-42'"
dbtool tx-search "send_packet.packet_connection='connection-41'"`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			txSearch(args[0])
		},
	}
	txSearchCmd.Flags().StringVar(&outfile, "out", outfile, "Output file")
	txSearchCmd.Flags().BoolVar(&full, "full", full, "Download all pages (alternatively, use --pages)")
	txSearchCmd.Flags().IntVar(&pages, "pages", pages, "Number of pages to download (--use full for all pages)")
	txSearchCmd.Flags().IntVar(&perPage, "per-page", perPage, "Number of results per page")
	return txSearchCmd
}

func txSearch(query string) {
	fmt.Printf("data directory: %s\n", datadir)
	fmt.Printf("output file: %s\n", outfile)
	fmt.Printf("query: %s\n", query)
	fmt.Printf("pages: %d\n", pages)
	fmt.Printf("per-page: %d\n", perPage)

	err := utils.OpenDB(datadir)
	if err != nil {
		panic(err)
	}

	results, err := utils.DoTxSearch(query, full, pages, perPage)
	if err != nil {
		panic(err)
	}

	events := utils.FilterEvents(
		results,
		func(_ string) bool { return true },
	)

	err = printEvents(events)
	if err != nil {
		panic(err)
	}
}

func printEvents(events []*utils.EventInfo) error {
	f, err := openOutputFile(outfile)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf("Writing events to %s...\n", outfile)

	datawriter := bufio.NewWriter(f)
	for _, info := range events {
		fmt.Fprintf(f, "tx: %s\n", info.TxHash)
		fmt.Fprintf(f, "type: %s\n", info.Type)
		fmt.Fprintf(f, "attributes:\n")
		for _, attr := range info.Attributes {
			fmt.Fprintf(f, "	%s\n", attr)
		}
		fmt.Fprintf(f, "\n")
	}
	datawriter.Flush()

	return nil
}

func openOutputFile(filename string) (*os.File, error) {
	os.Remove(filename)
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
}
