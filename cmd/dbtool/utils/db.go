package utils

import (
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/state/txindex/kv"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

func OpenDB(dataPath string) error {
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
