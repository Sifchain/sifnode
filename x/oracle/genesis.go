package oracle

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	if len(strings.TrimSpace(data.AdminAddress)) != 0 {
		adminAddress, err := sdk.AccAddressFromBech32(data.AdminAddress)
		if err != nil {
			panic(err)
		}
		keeper.SetAdminAccount(ctx, adminAddress)
	}

	for _, list := range data.ValidatorWhitelist {
		powers := list.ValidatorWhitelist
		for _, power := range powers.ValidatorPower {
			keeper.UpdateOracleWhiteList(ctx, list.NetworkDescriptor, power.ValidatorAddress, power.VotingPower)

		}
	}

	for _, prophecy := range data.Prophecies {
		keeper.SetProphecy(ctx, *prophecy)
	}

	for _, fee := range data.CrossChainFee {
		networkIdentity := types.NetworkIdentity{NetworkDescriptor: fee.NetworkDescriptor}
		keeper.SetCrossChainFeeObj(ctx, networkIdentity, fee.CrossChainFee)
	}

	for _, consensusNeeded := range data.ConsensusNeeded {
		networkIdentity := types.NetworkIdentity{NetworkDescriptor: consensusNeeded.NetworkDescriptor}
		keeper.SetConsensusNeeded(ctx, networkIdentity, *consensusNeeded.ConsensusNeeded)
	}

	for _, lockBurnSequence := range data.WitnessLockBurnSequence {
		keeper.SetWitnessLockBurnNonceObj(ctx, *lockBurnSequence.WitnessLockBurnSequenceKey, *lockBurnSequence.WitnessLockBurnSequence)
	}

	for _, prophecyInfo := range data.ProphecyInfo {
		keeper.SetProphecyInfoObj(ctx, prophecyInfo.ProphecyInfo)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {

	adminAcc := keeper.GetAdminAccount(ctx)
	whiteList := keeper.GetAllWhiteList(ctx)
	prophecies := keeper.GetProphecies(ctx)
	crossChainFee := keeper.GetAllCrossChainFeeConfig(ctx)
	consensusNeeded := keeper.GetAllConsensusNeeded(ctx)
	witnessLockBurnSequence := keeper.GetAllWitnessLockBurnSequence(ctx)
	prophecyInfo := keeper.GetAllProphecyInfo(ctx)

	dbProphecies := make([]*types.Prophecy, len(prophecies))
	for i := range prophecies {
		dbProphecies[i] = &prophecies[i]
	}
	return &types.GenesisState{
		ValidatorWhitelist:      whiteList,
		AdminAddress:            adminAcc.String(),
		Prophecies:              dbProphecies,
		CrossChainFee:           crossChainFee,
		ConsensusNeeded:         consensusNeeded,
		WitnessLockBurnSequence: witnessLockBurnSequence,
		ProphecyInfo:            prophecyInfo,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(state *types.GenesisState) error {
	for _, crossChainFee := range state.CrossChainFee {
		if !crossChainFee.CrossChainFee.IsValid() {
			return errors.New("crossChainFee is not valid")
		}
	}

	for _, consensusNeeded := range state.ConsensusNeeded {
		if consensusNeeded.ConsensusNeeded.ConsensusNeeded > 100 {
			return errors.New("consensusNeeded stored is too large")
		}
	}

	return nil
}
