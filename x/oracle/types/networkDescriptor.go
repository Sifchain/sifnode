package types

import (
	"bytes"
	"encoding/binary"
)

// TODO if we should pre-define the network ID

// MaxNetworkDescriptor record the max networkID, the ID should be an incremental value
var MaxNetworkDescriptor uint32 = 0

// NetworkDescriptor define the different network like Ethereum, Binance
type NetworkDescriptor struct {
	NetworkID uint32 `json:"network_id"`
}

// NewNetworkDescriptor get a new NetworkDescriptor instance
func NewNetworkDescriptor(networkID uint32) NetworkDescriptor {
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
	_ = binary.Write(bytebuf, binary.BigEndian, n.NetworkID)
	return append(WhiteListValidatorPrefix, bytebuf.Bytes()...)
}
