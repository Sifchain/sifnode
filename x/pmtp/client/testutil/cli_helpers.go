package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	clpcli "github.com/Sifchain/sifnode/x/clp/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
)

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
