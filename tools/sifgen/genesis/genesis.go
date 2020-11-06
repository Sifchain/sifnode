package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
)

func ReplaceStakingBondDenom(nodeHomeDir string) error {
	var genesis common.Genesis

	genesisPath := fmt.Sprintf("%s/config/%s", nodeHomeDir, utils.GenesisFile)

	body, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &genesis); err != nil {
		return err
	}

	genesis.AppState.Staking.Params.BondDenom = common.StakeTokenDenom
	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(genesisPath, content, 0600); err != nil {
		return err
	}

	return nil
}
