package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	clpcli "github.com/Sifchain/sifnode/x/clp/client/cli"
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

func QueryClpPoolsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{fmt.Sprintf("--%s=json", cli.OutputFlag)}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdPools(""), args)
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
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from.String()),
		fmt.Sprintf("--%s=%s", clpcli.FlagPmtpPeriodStartBlock, pmtpStart.String()),
		fmt.Sprintf("--%s=%s", clpcli.FlagPmtpPeriodEndBlock, pmtpEnd.String()),
		fmt.Sprintf("--%s=%s", clpcli.FlagPmtpPeriodEpochLength, epochLength.String()),
		fmt.Sprintf("--%s=%s", clpcli.FlagPeriodGovernanceRate, rGov.String()),
	}
	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	fmt.Print("args", args)

	return clitestutil.ExecTestCLICmd(clientCtx, clpcli.GetCmdUpdatePmtpParams(), args)
}
