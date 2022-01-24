package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	keeper "github.com/Sifchain/sifnode/x/margin/keeper"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
)

/* Query margin params from chain. */

func QueryMarginParams(clientCtx client.Context) (*margintypes.Params, error) {
	marginTypes := &margintypes.Params{
		LeverageMax:          keeper.GetLeverageParam(clientCtx),
		InterestRateMax:      keeper.GetInterestRateMax(clientCtx),
		InterestRateMin:      keeper.GetInterestRateMin(clientCtx),
		InterestRateIncrease: keeper.GetInterestRateIncrease(clientCtx),
		InterestRateDecrease: keeper.GetInterestRateDecrease(clientCtx),
		HealthGainFactor:     keeper.GetHealthGainFactor(clientCtx),
		EpochLength:          keeper.GetEpochLength,
	}
	return marginTypes, nil
}

/* Query pool data from chain. */

func QueryPools(clientCtx client.Context) ([]*clptypes.Pool, error) {
	pools := keeper.GetEnabledPools(clientCtx)
	return pools, nil
}

type Event struct {
	ID         int64
	EventType  string
	Height     int32
	Attributes []Attribute
	Metadata   string
}

type Attribute struct {
	Key   string
	Value string
}

/* Returns events from block_results?height=height */

func BlockEvents(clientCtx client.Context, height int64, baseUrl string) ([]*Event, error) {

	path := fmt.Sprintf("%s/block_results?height=%d", baseUrl, height)

	resp, err := http.Get(path)
	if err != nil {
		return nil, err.Error()
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err.Error()
	}

	var events []*Event
	json.Unmarshal(body, &events)

	return events, nil
}
