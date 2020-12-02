package types

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Events store all events in a block
type Events struct {
	BlockHash common.Hash
	Events    []EthereumEvent
}

// AddEvents append event
func (e Events) AddEvents(event EthereumEvent) (Events, error) {
	for _, n := range e.Events {
		if n.Equal(event) {
			return Events{}, errors.New("event already in list")
		}
	}
	return Events{
		BlockHash: e.BlockHash,
		Events:    append(e.Events, event),
	}, nil
}

// EthEventBuffer store all events from Ethereum smart contract
type EthEventBuffer struct {
	Buffer map[*big.Int]Events
}

// NewEthEventBuffer create a new instance of EthEventBuffer
func NewEthEventBuffer() EthEventBuffer {
	return EthEventBuffer{
		Buffer: make(map[*big.Int]Events),
	}
}

// AddEvent insert a new event to queue
func (buff EthEventBuffer) AddEvent(blockNumber *big.Int, blockHash common.Hash, event EthereumEvent) error {
	events, ok := buff.Buffer[blockNumber]
	if ok {
		if blockHash == events.BlockHash {
			newEvents, err := events.AddEvents(event)
			if err != nil {
				return err
			}
			buff.Buffer[blockNumber] = newEvents
		} else {
			// different hash with the same height
			return errors.New("different event's block hash with the same block height")
		}
	} else {
		buff.Buffer[blockNumber] = Events{
			BlockHash: blockHash,
			Events:    []EthereumEvent{event},
		}
	}
	return nil
}
