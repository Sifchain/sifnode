package utils

import (
	"encoding/hex"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type EventInfo struct {
	Type           string
	PacketSequence int
	ChannelID      string
	TxHash         string
	Denom          string
	Amount         uint64
	Attributes     []string
	RealAttributes []abcitypes.EventAttribute
}

func FilterEvents(
	txs []*coretypes.ResultTx,
	typeFilter func(string) bool,
) []*EventInfo {
	infos := []*EventInfo{}
	for _, tx := range txs {
		txHash := hex.EncodeToString(tx.Tx.Hash())
		for _, ev := range tx.TxResult.Events {
			if typeFilter(ev.Type) {
				attributes := []string{}
				for _, attr := range ev.Attributes {
					attributes = append(attributes, attr.String())
				}
				info := &EventInfo{
					TxHash:         txHash,
					Type:           ev.Type,
					Attributes:     attributes,
					RealAttributes: ev.Attributes,
				}
				infos = append(infos, info)
			}
		}
	}
	return infos
}
