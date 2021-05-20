package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/common/types"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
)

func ReplaceStakingBondDenom(nodeHomeDir string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.Staking.Params.BondDenom = common.StakeTokenDenom
	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceCLPMinCreatePoolThreshold(nodeHomeDir, minCreatePoolThreshold string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.CLP.Params.MinCreatePoolThreshold = minCreatePoolThreshold
	content, err := json.Marshal(genesis)
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
	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceGovDepositParamsMaxDepositPeriod(nodeHomeDir, period string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.Gov.DepositParams.MaxDepositPeriod = period
	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func ReplaceGovVotingParamsVotingPeriod(nodeHomeDir, period string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	(*genesis).AppState.Gov.VotingParams.VotingPeriod = period
	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func InitializeCLP(nodeHomeDir, clpConfigURL string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	response, err := http.Get(fmt.Sprintf("%v", clpConfigURL))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var pools types.CLP
	if err := json.Unmarshal(body, &pools); err != nil {
		return err
	}

	(*genesis).AppState.CLP.PoolList = pools.PoolList
	(*genesis).AppState.CLP.LiquidityProviderList = pools.LiquidityProviderList
	(*genesis).AppState.CLP.CLPModuleAddress = pools.CLPModuleAddress

	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}

func InitializeEthbridge(nodeHomeDir, ethbridgeConfigURL string) error {
	genesis, err := readGenesis(nodeHomeDir)
	if err != nil {
		return err
	}

	response, err := http.Get(fmt.Sprintf("%v", ethbridgeConfigURL))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var tokens types.Ethbridge
	if err := json.Unmarshal(body, &tokens); err != nil {
		return err
	}

	(*genesis).AppState.Ethbridge.PeggyTokens = tokens.PeggyTokens
	(*genesis).AppState.Ethbridge.CethReceiverAccount = tokens.CethReceiverAccount

	content, err := json.Marshal(genesis)
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
