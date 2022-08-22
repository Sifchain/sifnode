//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package node

import (
	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/genesis"
)

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_seedGenesis() error {
	if err := genesis.ReplaceMarginGenesis(common.DefaultNodeHome); err != nil {
		return err
	}
	return nil
}
