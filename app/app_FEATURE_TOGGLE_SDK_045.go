//go:build FEATURE_TOGGLE_SDK_045
// +build FEATURE_TOGGLE_SDK_045

package app

import (
	disptypes "github.com/Sifchain/sifnode/x/dispensation/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	ibctransferoverride "github.com/Sifchain/sifnode/x/ibctransfer"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
)

func FEATURE_TOGGLE_SDK_045_getOrderBeginBlockers(transferModule *ibctransferoverride.AppModule) []string {
	return []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		disptypes.ModuleName,
		transferModule.Name(),
		ethbridgetypes.ModuleName,
		tokenregistrytypes.ModuleName,
		oracletypes.ModuleName,
	}
}

func FEATURE_TOGGLE_SDK_045_getOrderEndBlockers(transferModule *ibctransferoverride.AppModule) []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ibchost.ModuleName,
		disptypes.ModuleName,
		transferModule.Name(),
		ethbridgetypes.ModuleName,
		tokenregistrytypes.ModuleName,
		oracletypes.ModuleName,
	}
}

// NOTE: Capability module must occur first so that it can initialize any capabilities
// so that other modules that want to create or claim capabilities afterwards in InitChain
// can do so safely.
func FEATURE_TOGGLE_SDK_045_getOrderInitGenesis(transferModule *ibctransferoverride.AppModule) []string {
	return []string{
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		disptypes.ModuleName,
		transferModule.Name(),
		ethbridgetypes.ModuleName,
	}
}
