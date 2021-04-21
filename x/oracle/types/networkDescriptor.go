package types

import (
	"bytes"
	"encoding/binary"
)

// MaxNetworkDescriptor record the max networkID, the ID should be an incremental value
var MaxNetworkDescriptor int32 = 0

// NetworkDescriptor define the different network like Ethereum, Binance
type NetworkDescriptor struct {
	NetworkID int32 `json:"network_id"`
}

// NewNetworkDescriptor get a new NetworkDescriptor instance
func NewNetworkDescriptor(networkID int32) NetworkDescriptor {
	if networkID > MaxNetworkDescriptor {
		MaxNetworkDescriptor = networkID
	}
	return NetworkDescriptor{
		NetworkID: networkID,
	}
}

// GetPrefix return storage prefix
func (n NetworkDescriptor) GetPrefix() []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, n.NetworkID)
	return append(WhiteListValidatorPrefix, bytebuf.Bytes()...)
}
