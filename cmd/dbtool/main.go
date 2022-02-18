package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/state/txindex/kv"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

var (
	datadir string
	outfile string
	query   string
	page    int
	perPage int
)

var (
	// "message.action='/ibc.core.client.v1.MsgUpdateClient'"
	// "fungible_token_packet.denom='ujuno'" // XXX what it the denom of uluna ibc/lkjljlkjlkj
	// "fungible_token_packet.sucess='001'"
	// "fungible_token_packet.denom='transfer/channel-19/ungm'"
	// "fungible_token_packet.denom='transfer/channel-18/uluna'"
	defaultQuery   = "update_client.client_id='07-tendermint-41'"
	defaultPage    = 1
	defaultPerPage = 10
)

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultDataDir := fmt.Sprintf("%s/.sifnoded/data", homedir)
	defaultOutFile := fmt.Sprintf("%s/dbtool.data", homedir)
	flag.StringVar(&datadir, "data", defaultDataDir, "data directory")
	flag.StringVar(&outfile, "out", defaultOutFile, "output file")
	flag.StringVar(&query, "query", defaultQuery, "query string")
	flag.IntVar(&page, "page", defaultPage, "page number")
	flag.IntVar(&perPage, "per-page", defaultPerPage, "results per page")
	flag.Parse()
	fmt.Printf("data directory: %s\n", datadir)
	fmt.Printf("output file: %s\n", outfile)
	fmt.Printf("query: %s\n", query)
	fmt.Printf("page: %d\n", page)
	fmt.Printf("per-page: %d\n", perPage)
}

func main() {

	err := openDB(datadir)
	if err != nil {
		panic(err)
	}

	f, err := openOutputFile(outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

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

	fmt.Printf("results: %d | total: %d\n", len(res.Txs), res.TotalCount)

	fmt.Printf("Writing transactions to %s...\n", outfile)
	datawriter := bufio.NewWriter(f)
	for _, tx := range res.Txs {
		for _, ev := range tx.TxResult.Events {
			datawriter.WriteString(ev.String() + "\n\n")
		}
		datawriter.WriteString("*************************************\n\n")
	}
	datawriter.Flush()
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
