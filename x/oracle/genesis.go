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
	if data.AddressWhitelist != nil {
		for networkDescriptor, list := range data.AddressWhitelist {
			keeper.SetOracleWhiteList(ctx, types.NewNetworkIdentity(types.NetworkDescriptor(networkDescriptor)), *list)
		}
	}

	if len(strings.TrimSpace(data.AdminAddress)) != 0 {
		adminAddress, err := sdk.AccAddressFromBech32(data.AdminAddress)
		if err != nil {
			panic(err)
		}
		keeper.SetAdminAccount(ctx, adminAddress)
	}

	for _, prophecy := range data.Prophecies {
		keeper.SetProphecy(ctx, *prophecy)
	}

	for key, value := range data.CrossChainFee {
		keeper.SetCrossChainFeeObj(ctx, types.NewNetworkIdentity(types.NetworkDescriptor(key)), value)
	}

	for key, value := range data.ConsensusNeeded {
		keeper.SetConsensusNeeded(ctx, types.NewNetworkIdentity(types.NetworkDescriptor(key)), value)
	}

	for key, value := range data.WitnessLockBurnSequence {
		keeper.SetWitnessLockBurnNonceViaRawKey(ctx, []byte(key), value)
	}

	for _, value := range data.ProphecyInfo {
		keeper.SetProphecyInfoObj(ctx, value)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	whiteList := keeper.GetAllWhiteList(ctx)
	wl := make(map[uint32]*types.ValidatorWhiteList, len(whiteList))
	for key, entry := range whiteList {
		wlEntry := entry
		wl[uint32(key)] = &wlEntry
	}
	adminAcc := keeper.GetAdminAccount(ctx)
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
		AddressWhitelist:        wl,
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
		if !crossChainFee.IsValid() {
			return errors.New("crossChainFee is not valid")
		}
	}

	for _, consensusNeeded := range state.ConsensusNeeded {
		if consensusNeeded > 100 {
			return errors.New("consensusNeeded stored is too large")
		}
	}

	return nil
}
