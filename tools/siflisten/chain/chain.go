package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	events "github.com/Sifchain/sifnode/tools/siflisten/events"
	"github.com/Sifchain/sifnode/x/clp/types"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsProposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

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
	fmt.Println("QUERY MARGIN PARAMS")

	var k paramsKeeper.Keeper
	var t paramsProposal.QueryParamsRequest

	resp, err := k.Params(context.Background(), &t)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(resp)

	//paramsKeeper.NewQuerier()

	//.Params(context.Background())
	//paramsKeeper.NewQuerier(sifapp.MarginKeeper, )
	//queryClient.GetMTP()
	// marginKeeper := sifapp.MarginKeeper

	// genesis := marginKeeper.ExportGenesis(ctx)

	// marginParams := margintypes.Params{
	// 	LeverageMax:          genesis.Params.LeverageMax,
	// 	InterestRateMax:      genesis.Params.InterestRateMax,
	// 	InterestRateMin:      genesis.Params.InterestRateMin,
	// 	InterestRateIncrease: genesis.Params.InterestRateIncrease,
	// 	InterestRateDecrease: genesis.Params.InterestRateDecrease,
	// 	HealthGainFactor:     genesis.Params.HealthGainFactor,
	// 	EpochLength:          genesis.Params.EpochLength,
	// }

	return nil, nil
}

/* Returns events from block_results?height=height */

func BlockEvents(clientCtx client.Context) ([]*events.Event, error) {
	//take a look on the return of this path and return events.Event
	fmt.Println("BlockEvents")
	fmt.Println(clientCtx.NodeURI)
	path := fmt.Sprintf("%s/block_results?height=%d", clientCtx.NodeURI, clientCtx.Height)

	resp, err := http.Get(path)

	if err != nil {
		fmt.Println("err1")
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err2")
		fmt.Println(err)
		return nil, err
	}

	json.Unmarshal(body, &blocksResults)
	fmt.Println("BlockEvents")
	fmt.Println(blocksResults)

	return nil, nil
}

func QueryPools(clientCtx client.Context) ([]*types.Pool, error) {
	queryClient := types.NewQueryClient(clientCtx)

	result, err := queryClient.GetPools(context.Background(), &types.PoolsReq{})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("QueryPools")
	fmt.Println(result)
	clientCtx.PrintProto(result)
	return result.Pools, nil
}
