package types

// TODO: This should be moved to new 'events' directory and expanded so that it can
// serve as a local store of witnessed events and allow for re-trying failed relays.

// EventRecords map of transaction hashes to EthereumEvent structs
var EventRecords = make(map[string]EthereumEvent)

// NewEventWrite add a validator's address to the official claims list
func NewEventWrite(txHash string, event EthereumEvent) {
	EventRecords[txHash] = event
}
