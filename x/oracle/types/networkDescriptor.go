package types

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// NetworkDescriptor define the different network like Ethereum, Binance
type NetworkDescriptor struct {
	NetworkID uint32 `json:"network_id"`
}

// NewNetworkDescriptor get a new NetworkDescriptor instance
func NewNetworkDescriptor(networkID uint32) NetworkDescriptor {
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

// GetFromPrefix return a NetworkDescriptor from prefix
func GetFromPrefix(key []byte) (NetworkDescriptor, error) {
	if len(key) == 5 {
		var data uint32
		bytebuff := bytes.NewBuffer(key[1:])
		binary.Read(bytebuff, binary.BigEndian, &data)
		return NewNetworkDescriptor(data), nil
	}

	return NetworkDescriptor{}, errors.New("prefix is invalid")
}
