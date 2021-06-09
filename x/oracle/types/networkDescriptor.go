package types

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// NetworkDescriptor define the different network like Ethereum, Binance
type NetworkDescriptor struct {
	NetworkID NetworkID `json:"network_id"`
}

// NewNetworkDescriptor get a new NetworkDescriptor instance
func NewNetworkDescriptor(networkID NetworkID) NetworkDescriptor {
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
		var data NetworkID
		bytebuff := bytes.NewBuffer(key[1:])
		err := binary.Read(bytebuff, binary.BigEndian, &data)
		if err == nil {
			return NewNetworkDescriptor(data), nil
		}
		return NetworkDescriptor{}, err
	}

	return NetworkDescriptor{}, errors.New("prefix is invalid")
}

// IsValid check if the network id is valid
func (n NetworkID) IsValid() bool {

	if n == NetworkID_NETWORK_ID_UNSPECIFIED {
		return false
	}

	_, ok := NetworkID_name[int32(n)]
	return ok
}
