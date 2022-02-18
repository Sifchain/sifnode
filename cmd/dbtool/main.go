package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/rpc/core"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
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
	// defaultQuery = "update_client.client_id='07-tendermint-41' AND fungible_token_packet.success='false'"
	defaultQuery = "update_client.client_id='07-tendermint-41'"
	// defaultQuery   = "send_packet.packet_dst_port='transfer'"
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

	events := filterEvents(res.Txs)

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

type EventInfo struct {
	Type           string
	PacketSequence int
	ChannelID      string
	TxHash         string
	Denom          string
	Amount         uint64
	Attributes     []string
}

func filterEvents(txs []*coretypes.ResultTx) []*EventInfo {
	infos := []*EventInfo{}
	for _, tx := range txs {
		txHash := hex.EncodeToString(tx.Tx.Hash())
		for _, ev := range tx.TxResult.Events {
			attributes := []string{}
			if keepEvent(ev.Type) {
				for _, attr := range ev.Attributes {
					attributes = append(attributes, attr.String())
				}
				info := &EventInfo{
					TxHash:     txHash,
					Type:       ev.Type,
					Attributes: attributes,
				}
				infos = append(infos, info)
			}
		}
	}
	return infos
}

func keepEvent(eventType string) bool {

	// incoming transfers
	// return eventType == "recv_packet" ||
	// 	eventType == "write_acknowledgement"

	// outgoing transfers
	return eventType == "send_packet" ||
		eventType == "acknowledge_packet" ||
		eventType == "timeout_packet"

	// return true

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
