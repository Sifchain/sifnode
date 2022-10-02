//go:build TEST_INTEGRATION
// +build TEST_INTEGRATION

package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	clpcli "github.com/Sifchain/sifnode/x/clp/client/cli"
	margincli "github.com/Sifchain/sifnode/x/margin/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/viper"
)

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

func QueryBalancesExec(clientCtx client.Context, address fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{address.String(), fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, bankcli.GetBalancesCmd(), args)
}

func QueryMarginParamsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, margincli.GetCmdParams(), args)
}

func QueryMarginPositionsForAddressExec(clientCtx client.Context, address fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{address.String(), fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, margincli.GetCmdQueryPositionsForAddress(), args)
}

func QueryClpPoolsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdPools(""), args)
}

func QueryClpPoolExec(clientCtx client.Context, symbol string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{symbol, fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdPool(""), args)
}

func QueryClpPmtpParamsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdPmtpParams(""), args)
}

func MsgClpUpdatePmtpParamsExec(clientCtx client.Context, from, pmtpStart, pmtpEnd, epochLength, rGov fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	viper.Set(flags.FlagFrom, from.String())
	viper.Set(clpcli.FlagPmtpPeriodStartBlock, pmtpStart.String())
	viper.Set(clpcli.FlagPmtpPeriodEndBlock, pmtpEnd.String())
	viper.Set(clpcli.FlagPmtpPeriodEpochLength, epochLength.String())
	viper.Set(clpcli.FlagPeriodGovernanceRate, rGov.String())

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, viper.Get(flags.FlagFrom)),
		fmt.Sprintf("--%s=%s", clpcli.FlagPmtpPeriodStartBlock, viper.Get(clpcli.FlagPmtpPeriodStartBlock)),
		fmt.Sprintf("--%s=%s", clpcli.FlagPmtpPeriodEndBlock, viper.Get(clpcli.FlagPmtpPeriodEndBlock)),
		fmt.Sprintf("--%s=%s", clpcli.FlagPmtpPeriodEpochLength, viper.Get(clpcli.FlagPmtpPeriodEpochLength)),
		fmt.Sprintf("--%s=%s", clpcli.FlagPeriodGovernanceRate, viper.Get(clpcli.FlagPeriodGovernanceRate)),
	}
	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdUpdatePmtpParams(), args)
}

func MsgClpEndPolicyExec(clientCtx client.Context, from fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	viper.Set(flags.FlagFrom, from.String())
	viper.Set(clpcli.FlagEndCurrentPolicy, "true")

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, viper.Get(flags.FlagFrom)),
		fmt.Sprintf("--%s=%s", clpcli.FlagEndCurrentPolicy, viper.Get(clpcli.FlagEndCurrentPolicy)),
	}
	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdModifyPmtpRates(), args)
}

func MsgClpModifyPmtpRatesExec(clientCtx client.Context, from, blockRate, runningRate fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	viper.Set(flags.FlagFrom, from.String())
	viper.Set(clpcli.FlagBlockRate, blockRate.String())
	viper.Set(clpcli.FlagRunningRate, runningRate.String())

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, viper.Get(flags.FlagFrom)),
		fmt.Sprintf("--%s=%s", clpcli.FlagBlockRate, viper.Get(clpcli.FlagBlockRate)),
		fmt.Sprintf("--%s=%s", clpcli.FlagRunningRate, viper.Get(clpcli.FlagRunningRate)),
	}
	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdModifyPmtpRates(), args)
}

func MsgMarginOpenExec(clientCtx client.Context, from fmt.Stringer, collateralAsset string, collateralAmount fmt.Stringer, borrowAsset, position string, leverage fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	viper.Set(flags.FlagFrom, from.String())
	viper.Set("collateral_asset", collateralAsset)
	viper.Set("collateral_amount", collateralAmount.String())
	viper.Set("borrow_asset", borrowAsset)
	viper.Set("position", position)
	viper.Set("leverage", leverage)

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, viper.Get(flags.FlagFrom)),
		fmt.Sprintf("--%s=%s", "collateral_asset", viper.Get("collateral_asset")),
		fmt.Sprintf("--%s=%s", "collateral_amount", viper.Get("collateral_amount")),
		fmt.Sprintf("--%s=%s", "borrow_asset", viper.Get("borrow_asset")),
		fmt.Sprintf("--%s=%s", "position", viper.Get("position")),
		fmt.Sprintf("--%s=%s", "leverage", viper.Get("leverage")),
	}
	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, margincli.GetOpenCmd(), args)
}

func MsgMarginCloseExec(clientCtx client.Context, from fmt.Stringer, id uint64, extraArgs ...string) (testutil.BufferWriter, error) {
	viper.Set(flags.FlagFrom, from.String())
	viper.Set("id", id)

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, viper.Get(flags.FlagFrom)),
		fmt.Sprintf("--%s=%d", "id", viper.Get("id")),
	}
	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, margincli.GetCloseCmd(), args)
}
