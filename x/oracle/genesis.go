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

	for _, list := range data.NetworkConfigData {
		powers := list.ValidatorWhitelist
		if powers != nil {
			for _, power := range powers.ValidatorPower {
				err := keeper.UpdateOracleWhiteList(ctx, list.NetworkDescriptor, power.ValidatorAddress, power.VotingPower)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	for _, prophecy := range data.Prophecies {
		keeper.SetProphecy(ctx, *prophecy)
	}

	for _, fee := range data.NetworkConfigData {
		if fee != nil && fee.CrossChainFee != nil {
			networkIdentity := types.NetworkIdentity{NetworkDescriptor: fee.NetworkDescriptor}
			keeper.SetCrossChainFeeObj(ctx, networkIdentity, fee.CrossChainFee)
		}
	}

	for _, consensusNeeded := range data.NetworkConfigData {
		if consensusNeeded != nil && consensusNeeded.ConsensusNeeded != nil {
			networkIdentity := types.NetworkIdentity{NetworkDescriptor: consensusNeeded.NetworkDescriptor}
			keeper.SetConsensusNeeded(ctx, networkIdentity, *consensusNeeded.ConsensusNeeded)
		}
	}

	for _, lockBurnSequence := range data.WitnessLockBurnSequence {
		keeper.SetWitnessLockBurnSequenceObj(ctx, *lockBurnSequence.WitnessLockBurnSequenceKey, *lockBurnSequence.WitnessLockBurnSequence)
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

	// get the init data from whitelist
	networkConfigData := whiteList

	// merge the data from crossChainFee
	for _, fee := range crossChainFee {
		found := false
		for index := 0; index < len(networkConfigData); index++ {
			if fee.NetworkDescriptor == networkConfigData[index].NetworkDescriptor {
				networkConfigData[index].CrossChainFee = fee.CrossChainFee
				found = true
				break
			}
		}
		if !found {
			networkConfigData = append(networkConfigData, &types.NetworkConfigData{
				NetworkDescriptor: fee.NetworkDescriptor,
				CrossChainFee:     fee.CrossChainFee,
			})

		}
	}

	// merge the data from consensusNeeded
	for _, consensus := range consensusNeeded {
		found := false
		for index := 0; index < len(networkConfigData); index++ {
			if consensus.NetworkDescriptor == networkConfigData[index].NetworkDescriptor {
				networkConfigData[index].ConsensusNeeded = consensus.ConsensusNeeded
				found = true
				break
			}
		}
		if !found {
			networkConfigData = append(networkConfigData, &types.NetworkConfigData{
				NetworkDescriptor: consensus.NetworkDescriptor,
				ConsensusNeeded:   consensus.ConsensusNeeded,
			})

		}
	}

	return &types.GenesisState{
		NetworkConfigData:       networkConfigData,
		AdminAddress:            adminAcc.String(),
		Prophecies:              dbProphecies,
		WitnessLockBurnSequence: witnessLockBurnSequence,
		ProphecyInfo:            prophecyInfo,
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(state *types.GenesisState) error {
	for _, networkConfigData := range state.NetworkConfigData {
		if !networkConfigData.CrossChainFee.IsValid() {
			return errors.New("crossChainFee is not valid")
		}
	}

	for _, networkConfigData := range state.NetworkConfigData {
		if networkConfigData.ConsensusNeeded.ConsensusNeeded > 100 {
			return errors.New("consensusNeeded stored is too large")
		}
	}

	return nil
}
