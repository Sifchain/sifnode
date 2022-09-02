//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package genesis

import (
	"encoding/json"
	"fmt"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

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
	(*genesis).AppState.Margin.Params.ForceCloseThreshold = gen.Params.ForceCloseThreshold.String()
	(*genesis).AppState.Margin.Params.RemovalQueueThreshold = gen.Params.RemovalQueueThreshold.String()
	(*genesis).AppState.Margin.Params.MaxOpenPositions = json.Number(fmt.Sprintf("%d", gen.Params.MaxOpenPositions))
	(*genesis).AppState.Margin.Params.PoolOpenThreshold = gen.Params.PoolOpenThreshold.String()
	(*genesis).AppState.Margin.Params.ForceCloseFundPercentage = gen.Params.ForceCloseFundPercentage.String()
	(*genesis).AppState.Margin.Params.ForceCloseInsuranceFundAddress = gen.Params.ForceCloseInsuranceFundAddress
	(*genesis).AppState.Margin.Params.IncrementalInterestPaymentFundPercentage = gen.Params.IncrementalInterestPaymentFundPercentage.String()
	(*genesis).AppState.Margin.Params.IncrementalInterestPaymentInsuranceFundAddress = gen.Params.IncrementalInterestPaymentInsuranceFundAddress
	(*genesis).AppState.Margin.Params.SqModifier = gen.Params.SqModifier.String()
	(*genesis).AppState.Margin.Params.SafetyFactor = gen.Params.SafetyFactor.String()
	(*genesis).AppState.Margin.Params.ClosedPools = gen.Params.ClosedPools
	(*genesis).AppState.Margin.Params.IncrementalInterestPaymentEnabled = gen.Params.IncrementalInterestPaymentEnabled
	(*genesis).AppState.Margin.Params.WhitelistingEnabled = gen.Params.WhitelistingEnabled

	content, err := tmjson.Marshal(genesis)
	if err != nil {
		return err
	}

	if err := writeGenesis(nodeHomeDir, content); err != nil {
		return err
	}

	return nil
}
