//go:build !FEATURE_TOGGLE_SDK_045
// +build !FEATURE_TOGGLE_SDK_045

package app

import (
	ibctransferoverride "github.com/Sifchain/sifnode/x/ibctransfer"
)

func FEATURE_TOGGLE_SDK_045_getOrderBeginBlockers(transferModule *ibctransferoverride.AppModule) []string {
	return []string{}
}

func FEATURE_TOGGLE_SDK_045_getOrderEndBlockers(transferModule *ibctransferoverride.AppModule) []string {
	return []string{}
}

func FEATURE_TOGGLE_SDK_045_getOrderInitGenesis(transferModule *ibctransferoverride.AppModule) []string {
	return []string{}
}
