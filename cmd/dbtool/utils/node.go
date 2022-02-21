package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/Sifchain/sifnode/app"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/types"
)

func LoadGenesis(datadir string) (*types.GenesisDoc, error) {
	config := cfg.DefaultConfig()
	config.DBPath = datadir

	stateDB, err := GetStateDB(config)
	if err != nil {
		return nil, err
	}

	_, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(
		stateDB,
		node.DefaultGenesisDocProviderFunc(config),
	)
	if err != nil {
		return nil, err
	}

	return genDoc, nil
}

func NewSifApp(datadir string) (*app.SifchainApp, error) {
	encCfg := app.MakeTestEncodingConfig() // Ideally, we would reuse the one created by NewRootCmd.
	encCfg.Marshaler = codec.NewProtoCodec(encCfg.InterfaceRegistry)

	appDB, err := sdk.NewLevelDB("application", datadir)
	if err != nil {
		return nil, err
	}

	traceWriterFile := fmt.Sprintf("/tmp/ibc.trace")
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return nil, err
	}

	sifApp := app.NewSifApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		appDB,
		traceWriter,
		false,
		map[int64]bool{},
		"",
		0,
		encCfg,
		app.EmptyAppOptions{},
	)

	err = sifApp.LoadLatestVersion()
	if err != nil {
		return nil, err
	}

	return sifApp, nil
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666,
	)
}

// ************************************************************************
// Initialize siff app

// validators := make([]*types.Validator, len(genDoc.Validators))
// for i, val := range genDoc.Validators {
// 	validators[i] = types.NewValidator(val.PubKey, val.Power)
// }
// validatorSet := types.NewValidatorSet(validators)
// nextVals := types.TM2PB.ValidatorUpdates(validatorSet)
// csParams := types.TM2PB.ConsensusParams(genDoc.ConsensusParams)
// req := abci.RequestInitChain{
// 	Time:            genDoc.GenesisTime,
// 	ChainId:         genDoc.ChainID,
// 	InitialHeight:   genDoc.InitialHeight,
// 	ConsensusParams: csParams,
// 	Validators:      nextVals,
// 	AppStateBytes:   genDoc.AppState,
// }

// sifApp.InitChain(req)

// ************************************************************************
// Create context
