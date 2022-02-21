/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/rpc/core"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/state/txindex/kv"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
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
		Args:  cobra.MinimumNArgs(1),
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

	err := openDB(datadir)
	if err != nil {
		panic(err)
	}

	results := doQuery(query, full, pages, perPage)

	events := filterEvents(results)

	err = printEvents(events)
	if err != nil {
		panic(err)
	}
}

func openDB(dataPath string) error {
	config := cfg.DefaultConfig()
	config.DBPath = dataPath

	blockStore, err := getBlockstore(config)
	if err != nil {
		return err
	}

	txIndexer, err := getTxIndexer(config)
	if err != nil {
		return err
	}

	core.SetEnvironment(
		&core.Environment{
			BlockStore: blockStore,
			TxIndexer:  txIndexer,
		},
	)

	return nil
}

func getBlockstore(config *cfg.Config) (*store.BlockStore, error) {
	db, err := dbm.NewDB(
		"blockstore",
		dbm.BackendType(config.DBBackend),
		config.DBDir(),
	)
	if err != nil {
		return nil, err
	}
	return store.NewBlockStore(db), nil
}

func getTxIndexer(config *cfg.Config) (*kv.TxIndex, error) {
	db, err := dbm.NewDB(
		"tx_index",
		dbm.BackendType(config.DBBackend),
		config.DBDir(),
	)
	if err != nil {
		return nil, err
	}
	return kv.NewTxIndex(db), nil
}

func openOutputFile(filename string) (*os.File, error) {
	os.Remove(filename)
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
}

func doQuery(query string, full bool, pages, perPage int) []*coretypes.ResultTx {
	results := []*coretypes.ResultTx{}
	page := 1
	keepGoing := true
	for keepGoing {
		fmt.Printf("Getting transactions (page %d, perPage %d)...\n", page, perPage)
		res, err := core.TxSearch(
			&rpctypes.Context{},
			query,
			false,
			&page,
			&perPage,
			"asc",
		)
		if err != nil {
			panic(err)
		}
		results = append(results, res.Txs...)
		fmt.Printf("results: %d | total: %d\n", len(results), res.TotalCount)
		if len(res.Txs) < perPage ||
			(page == pages && !full) {
			keepGoing = false
		} else {
			page++
		}
	}
	return results
}

type EventInfo struct {
	Type           string
	PacketSequence int
	ChannelID      string
	TxHash         string
	Denom          string
	Amount         uint64
	Attributes     []string
	RealAttributes []abcitypes.EventAttribute
}

func filterEvents(txs []*coretypes.ResultTx) []*EventInfo {
	infos := []*EventInfo{}
	for _, tx := range txs {
		txHash := hex.EncodeToString(tx.Tx.Hash())
		for _, ev := range tx.TxResult.Events {
			if keepEvent(ev.Type) {
				attributes := []string{}
				for _, attr := range ev.Attributes {
					attributes = append(attributes, attr.String())
				}
				info := &EventInfo{
					TxHash:         txHash,
					Type:           ev.Type,
					Attributes:     attributes,
					RealAttributes: ev.Attributes,
				}
				infos = append(infos, info)
			}
		}
	}
	return infos
}

func keepEvent(eventType string) bool {
	// return true

	// incoming transfers
	// return eventType == "recv_packet" ||
	// 	eventType == "write_acknowledgement"

	// outgoing transfers
	return eventType == "send_packet" ||
		eventType == "acknowledge_packet" ||
		eventType == "timeout_packet"
}

func printEvents(events []*EventInfo) error {

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
