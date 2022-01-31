package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sifchain/sifnode/app"
	events "github.com/Sifchain/sifnode/tools/siflisten/events"
	"github.com/Sifchain/sifnode/x/clp/types"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
)

var sifapp app.SifchainApp

var blocksResults BlockResults

type Parts struct {
	Total int64  `json:"total"`
	Hash  string `json:"hash"`
}

type Version struct {
	Block string `json:"block"`
	App   string `json:"app"`
}

type BlockHeader struct {
	Version            Version `json:"version"`
	ChainID            string  `json:"chain_id"`
	Height             string  `json:"height"`
	Time               string  `json:"time"`
	LastBlockID        BlockID `json:"last_block_id"`
	LastCommitHash     string  `json:"last_commit_hash"`
	DataHash           string  `json:"data_hash"`
	ValidatorsHash     string  `json:"validators_hash"`
	NextValidatorsHash string  `json:"next_validators_hash"`
	ConsensusHash      string  `json:"consensus_hash"`
	LastResultsHash    string  `json:"last_results_hash"`
	AppHash            string  `json:"app_hash"`
	EvidenceHash       string  `json:"evidence_hash"`
	ProposerAddress    string  `json:"proposer_hash"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Validator struct {
	PubKey      PubKey
	VotingPower int64  `json:"voting_power"`
	Address     string `json:"address"`
}

type Evidence struct {
	Type             string    `json:"type"`
	Height           int64     `json:"height"`
	Time             int64     `json:"time"`
	TotalVotingPower int64     `json:"total_voting_power"`
	Validator        Validator `json:"validator"`
}

type Commit struct {
	Type             int64   `json:"type"`
	Height           string  `json:"height"`
	Round            int64   `json:"round"`
	BlockID          BlockID `json:"block_id"`
	Timestamp        string  `json:"timestamp"`
	ValidatorAddress string  `json:"validator_address"`
	ValidatorIndex   int64   `json:"validator_index"`
	Signature        string  `json:"signature"`
}

type LastCommit struct {
	Height     int64    `json:"height"`
	Round      int64    `json:"round"`
	BlockID    BlockID  `json:"block_id"`
	Signatures []Commit `json:"signatures"`
}

type Block struct {
	Header     BlockHeader `json:"header"`
	Data       []string    `json:"data"`
	Evidence   []Evidence  `json:"evidence"`
	LastCommit `json:"last_commit"`
}

type BlockID struct {
	Hash  string `json:"hash"`
	Parts Parts  `json:"parts"`
}

type BlockComplete struct {
	BlockID BlockID `json:"block_id"`
	Block   Block   `json:"block"`
}

type Result struct {
	Blocks      []BlockComplete `json:"blocks"`
	Total_count int64           `json:"total_count"`
}

type BlockResults struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
	Result  Result `json:"result"`
}

type ErrorResponse struct {
	Id      int64  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func QueryMarginParams(clientCtx client.Context) (*margintypes.Params, error) {

	queryClient := margintypes.NewQueryClient(clientCtx)
	queryClient.GetMTP()
	marginKeeper := sifapp.MarginKeeper

	genesis := marginKeeper.ExportGenesis(ctx)

	marginParams := margintypes.Params{
		LeverageMax:          genesis.Params.LeverageMax,
		InterestRateMax:      genesis.Params.InterestRateMax,
		InterestRateMin:      genesis.Params.InterestRateMin,
		InterestRateIncrease: genesis.Params.InterestRateIncrease,
		InterestRateDecrease: genesis.Params.InterestRateDecrease,
		HealthGainFactor:     genesis.Params.HealthGainFactor,
		EpochLength:          genesis.Params.EpochLength,
	}

	return &marginParams, nil
}

/* Returns events from block_results?height=height */

func BlockEvents(ctx sdk.Context, height int64) ([]*events.Event, error) {
	//take a look on the return of this path and return events.Event
	path := fmt.Sprintf("%s/block_results?height=%d", "https://rpc.cosmos.network", height)

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &blocksResults)

	return blocksResults, nil
}

func QueryPools(clientCtx client.Context, flagSet *pflag.FlagSet) (*clptypes.PoolsRes, error) {

	queryClient := types.NewQueryClient(clientCtx)

	pageReq, err := client.ReadPageRequest(flagSet)
	if err != nil {
		return nil, err
	}

	result, err := queryClient.GetPools(context.Background(), &types.PoolsReq{
		Pagination: pageReq,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
