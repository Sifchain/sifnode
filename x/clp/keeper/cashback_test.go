package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CalcCashbackAmount(t *testing.T) {
	rowanCashbacked := sdk.NewDec(10)
	totalPoolUnits := sdk.NewUint(999)
	lpPoolUnits := sdk.NewUint(333)
	expectedAmount := sdk.NewUint(3)

	amount := keeper.CalcCashbackAmount(rowanCashbacked, totalPoolUnits, lpPoolUnits)

	require.Equal(t, expectedAmount, amount)
}
