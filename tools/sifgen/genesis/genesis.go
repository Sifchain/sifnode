package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
)

func ReplaceStakingBondDenom(nodeHomeDir string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	genesis.AppState.Staking.Params.BondDenom = common.StakeTokenDenom
	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceMintBondDenom(nodeHomeDir string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	genesis.AppState.Mint.Params.MintDenom = common.StakeTokenDenom
	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceCLPMinCreatePoolThreshold(nodeHomeDir string, minCreatePoolThreshold uint64) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.CLP.Params.MinCreatePoolThreshold = json.Number(fmt.Sprintf("%d", minCreatePoolThreshold))
	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceGovDepositParamsMinDeposit(nodeHomeDir, tokenDenom string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.Gov.DepositParams.MinDeposit[0].Denom = tokenDenom
	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceGovDepositParamsMaxDepositPeriod(nodeHomeDir string, period time.Duration) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.Gov.DepositParams.MaxDepositPeriod = fmt.Sprintf("%v", period)
	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceGovVotingParamsVotingPeriod(nodeHomeDir string, period time.Duration) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.Gov.VotingParams.VotingPeriod = fmt.Sprintf("%v", period)
	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceMarginGenesis(nodeHomeDir string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	gen := margintypes.DefaultGenesis()
	(*genesis).AppState.Margin.Params.LeverageMax = gen.Params.LeverageMax.String()
	(*genesis).AppState.Margin.Params.InterestRateMax = gen.Params.InterestRateMax.String()
	(*genesis).AppState.Margin.Params.InterestRateMin = gen.Params.InterestRateMin.String()
	(*genesis).AppState.Margin.Params.InterestRateIncrease = gen.Params.InterestRateIncrease.String()
	(*genesis).AppState.Margin.Params.InterestRateDecrease = gen.Params.InterestRateDecrease.String()
	(*genesis).AppState.Margin.Params.HealthGainFactor = gen.Params.HealthGainFactor.String()
	(*genesis).AppState.Margin.Params.EpochLength = json.Number(fmt.Sprintf("%d", gen.Params.EpochLength))
	(*genesis).AppState.Margin.Params.Pools = gen.Params.Pools
	(*genesis).AppState.Margin.Params.RemovalQueueThreshold = gen.Params.RemovalQueueThreshold.String()
	(*genesis).AppState.Margin.Params.MaxOpenPositions = json.Number(fmt.Sprintf("%d", gen.Params.MaxOpenPositions))
	(*genesis).AppState.Margin.Params.PoolOpenThreshold = gen.Params.PoolOpenThreshold.String()
	(*genesis).AppState.Margin.Params.ForceCloseFundPercentage = gen.Params.ForceCloseFundPercentage.String()
	(*genesis).AppState.Margin.Params.ForceCloseFundAddress = gen.Params.ForceCloseFundAddress
	(*genesis).AppState.Margin.Params.IncrementalInterestPaymentFundPercentage = gen.Params.IncrementalInterestPaymentFundPercentage.String()
	(*genesis).AppState.Margin.Params.IncrementalInterestPaymentFundAddress = gen.Params.IncrementalInterestPaymentFundAddress
	(*genesis).AppState.Margin.Params.SqModifier = gen.Params.SqModifier.String()
	(*genesis).AppState.Margin.Params.SafetyFactor = gen.Params.SafetyFactor.String()
	(*genesis).AppState.Margin.Params.ClosedPools = gen.Params.ClosedPools
	(*genesis).AppState.Margin.Params.IncrementalInterestPaymentEnabled = gen.Params.IncrementalInterestPaymentEnabled
	(*genesis).AppState.Margin.Params.WhitelistingEnabled = gen.Params.WhitelistingEnabled
	(*genesis).AppState.Margin.Params.RowanCollateralEnabled = gen.Params.RowanCollateralEnabled

	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func readGenesis(nodeHomeDir string) (*common.Genesis, error) {
	var genesis common.Genesis

	genesisPath := fmt.Sprintf("%s/config/%s", nodeHomeDir, utils.GenesisFile)

	body, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &genesis); err != nil {
		return nil, err
	}

	return &genesis, nil
}

func writeGenesis(nodeHomeDir string, content []byte) error {
	genesisPath := fmt.Sprintf("%s/config/%s", nodeHomeDir, utils.GenesisFile)
	if err := ioutil.WriteFile(genesisPath, content, 0600); err != nil {
		return err
	}

	return nil
}
