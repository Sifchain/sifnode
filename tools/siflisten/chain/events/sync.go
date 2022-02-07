package events

import (
	"database/sql"
	"fmt"

	"github.com/Sifchain/sifnode/tools/siflisten/chain"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
)

func Sync(clientCtx sdkclient.Context, _ *sql.DB) {
	_, err := chain.QueryMarginParams(clientCtx)
	if err != nil {
		fmt.Println("err-querymarginparams")
		fmt.Println(err)
	}
	_, err = chain.QueryPools(clientCtx)
	if err != nil {
		fmt.Println("err-querypools")
		fmt.Println(err)
	}
	_, err = chain.BlockEvents(clientCtx)
	if err != nil {
		fmt.Println("err-blockevents")
		fmt.Println(err)
	}

}
