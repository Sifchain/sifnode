package types

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EventsInBlock store all events in a block, parent hash used to determine the best chain
type EventsInBlock struct {
	ParentHash common.Hash
	Events     []EthereumEvent
}

// NewEventsInBlock create new instance with parent hash
func NewEventsInBlock() EventsInBlock {
	return EventsInBlock{
		ParentHash: common.Hash{},
		Events:     []EthereumEvent{},
	}
}

// AddEvent append new event to list
func (e *EventsInBlock) AddEvent(event EthereumEvent) {
	// avoid add the same event twice
	for _, n := range e.Events {
		if n.Equal(event) {
			log.Println("EventsInBlock event already in list")
			return
		}
	}
	e.Events = append(e.Events, event)
}

// EventsInHeight store all events at the same height
type EventsInHeight struct {
	// map the block hash to its parent hash and event list
	EventsMap map[common.Hash]*EventsInBlock
}

// NewEventsInHeight create a new instance
func NewEventsInHeight() EventsInHeight {
	return EventsInHeight{
		EventsMap: make(map[common.Hash]*EventsInBlock),
	}
}

// AddEvent append event
func (e *EventsInHeight) AddEvent(blockHash common.Hash, event EthereumEvent) {
	events, ok := e.EventsMap[blockHash]
	if ok {
		events.AddEvent(event)
	} else {
		newEventsInBlock := NewEventsInBlock()
		newEventsInBlock.AddEvent(event)
		e.EventsMap[blockHash] = &newEventsInBlock
	}
}

// AddHeader add a new block hash into map
func (e *EventsInHeight) AddHeader(blockHash common.Hash, parentHash common.Hash) {
	events, ok := e.EventsMap[blockHash]
	// the events list the block hash already existed, then update the parent hash
	if ok {
		events.ParentHash = parentHash
	} else {
		newEventsInBlock := NewEventsInBlock()
		newEventsInBlock.ParentHash = parentHash
		e.EventsMap[blockHash] = &newEventsInBlock
	}
}

// EthEventBuffer store all events from Ethereum smart contract
type EthEventBuffer struct {
	Buffer    map[string]EventsInHeight
	MinHeight *big.Int
}

// NewEthEventBuffer create a new instance of EthEventBuffer
func NewEthEventBuffer() EthEventBuffer {
	return EthEventBuffer{
		Buffer:    make(map[string]EventsInHeight),
		MinHeight: big.NewInt(0),
	}
}

// AddEvent insert a new event to queue
// func (buff *EthEventBuffer) AddEvent(blockNumber *big.Int, blockHash common.Hash, event EthereumEvent) {
func (buff *EthEventBuffer) AddEvent(blockNumber fmt.Stringer, blockHash common.Hash, event EthereumEvent) {
	// Check if block number already in the map
	events, ok := buff.Buffer[blockNumber.String()]
	if ok {
		events.AddEvent(blockHash, event)
	} else {
		newEvents := NewEventsInHeight()
		newEvents.AddEvent(blockHash, event)
		buff.Buffer[blockNumber.String()] = newEvents
	}
}

// AddHeader create new entry for new header
func (buff *EthEventBuffer) AddHeader(blockNumber *big.Int, blockHash common.Hash, parentHash common.Hash) {
	if buff.MinHeight.Cmp(big.NewInt(0)) == 0 {
		buff.MinHeight = blockNumber
	}
	// Check if block number already in the map
	eventsInHeight, ok := buff.Buffer[blockNumber.String()]
	if ok {
		eventsInHeight.AddHeader(blockHash, parentHash)
	} else {
		newEventsInHeight := NewEventsInHeight()
		newEventsInHeight.AddHeader(blockHash, parentHash)
		buff.Buffer[blockNumber.String()] = newEventsInHeight
	}
}

// GetDepth get the depth of a block
func (buff *EthEventBuffer) GetDepth(blockNumber *big.Int, blockHash common.Hash) uint64 {
	eventsInHeight, ok := buff.Buffer[blockNumber.String()]
	if ok {
		// if there is block's parent is the block hash
		for key, eventsInBlock := range eventsInHeight.EventsMap {
			if eventsInBlock.ParentHash == blockHash {
				one := big.NewInt(1)
				one.Add(one, blockNumber)

				// recursive to its child block
				return buff.GetDepth(one, key) + 1
			}
		}
	}
	return 0
}

// RemoveHeight remove an entry
func (buff *EthEventBuffer) RemoveHeight() {
	delete(buff.Buffer, buff.MinHeight.String())
	buff.MinHeight.Add(buff.MinHeight, big.NewInt(1))
}

// GetHeaderEvents get the events in block of best chain
func (buff *EthEventBuffer) GetHeaderEvents() []EthereumEvent {
	eventsInHeight, ok := buff.Buffer[buff.MinHeight.String()]
	if ok {
		maxDepth := uint64(0)
		var result []EthereumEvent
		one := big.NewInt(1)
		one.Add(one, buff.MinHeight)
		for blockHash, eventsInBlock := range eventsInHeight.EventsMap {

			depth := buff.GetDepth(one, blockHash)

			if depth >= maxDepth {
				maxDepth = depth
				result = eventsInBlock.Events
			}
		}

		return result
	}

	return []EthereumEvent{}
}
