package clp

import (
	"math/rand"

	"github.com/Sifchain/sifnode/testutil/sample"
	clpsimulation "github.com/Sifchain/sifnode/x/clp/simulation"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = clpsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgAddLiquidityToRewardsBucket = "op_weight_msg_add_liquidity_to_rewards_bucket" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddLiquidityToRewardsBucket int = 100
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	clpGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&clpGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAddLiquidityToRewardsBucket int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddLiquidityToRewardsBucket, &weightMsgAddLiquidityToRewardsBucket, nil,
		func(_ *rand.Rand) {
			weightMsgAddLiquidityToRewardsBucket = defaultWeightMsgAddLiquidityToRewardsBucket
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddLiquidityToRewardsBucket,
		clpsimulation.SimulateMsgAddLiquidityToRewardsBucket(am.bankKeeper, am.keeper),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			opWeightMsgAddLiquidityToRewardsBucket,
			defaultWeightMsgAddLiquidityToRewardsBucket,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
				clpsimulation.SimulateMsgAddLiquidityToRewardsBucket(am.bankKeeper, am.keeper)
				return nil
			},
		),
	}
}
