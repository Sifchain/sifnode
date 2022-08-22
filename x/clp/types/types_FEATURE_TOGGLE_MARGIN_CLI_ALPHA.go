//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) ExtractDebt(X, Y sdk.Uint, toRowan bool) (sdk.Uint, sdk.Uint) {

	if toRowan {
		Y = Y.Add(p.NativeCustody).Add(p.NativeLiabilities)
		X = X.Add(p.ExternalCustody).Add(p.ExternalLiabilities)
	} else {
		X = X.Add(p.NativeCustody).Add(p.NativeLiabilities)
		Y = Y.Add(p.ExternalCustody).Add(p.ExternalLiabilities)
	}

	return X, Y
}
