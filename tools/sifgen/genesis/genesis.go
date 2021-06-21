package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
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

func readGenesis(nodeHomeDir string) (*common.Genesis, error) {
	var genesis common.Genesis

	genesisPath := fmt.Sprintf("%s/config/%s", nodeHomeDir, utils.GenesisFile)

	body, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	if err := tmjson.Unmarshal(body, &genesis); err != nil {
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
