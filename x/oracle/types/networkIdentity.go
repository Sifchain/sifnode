package types

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// NetworkIdentity define the different network like Ethereum, Binance
type NetworkIdentity struct {
	NetworkDescriptor NetworkDescriptor `json:"network_descriptor"`
}

// NewNetworkIdentity get a new NetworkIdentity instance
func NewNetworkIdentity(networkDescriptor NetworkDescriptor) NetworkIdentity {
	return NetworkIdentity{
		NetworkDescriptor: networkDescriptor,
	}
}

// GetPrefix return storage prefix
func (n NetworkIdentity) GetPrefix() []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, n.NetworkDescriptor)
	return append(WhiteListValidatorPrefix, bytebuf.Bytes()...)
}

// GetNativeTokenPrefix return storage prefix
func (n NetworkIdentity) GetNativeTokenPrefix() []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, n.NetworkDescriptor)
	return append(NativeTokenPrefix, bytebuf.Bytes()...)
}

// GetFromPrefix return a NetworkIdentity from prefix
func GetFromPrefix(key []byte) (NetworkIdentity, error) {
	if len(key) == 5 {
		var data NetworkDescriptor
		bytebuff := bytes.NewBuffer(key[1:])
		err := binary.Read(bytebuff, binary.BigEndian, &data)
		if err == nil {
			return NewNetworkIdentity(data), nil
		}
		return NetworkIdentity{}, err
	}

	return NetworkIdentity{}, errors.New("prefix is invalid")
}

// IsValid check if the network id is valid
func (n NetworkDescriptor) IsValid() bool {

	if n == NetworkDescriptor_NETWORK_DESCRIPTOR_UNSPECIFIED {
		return false
	}

	_, ok := NetworkDescriptor_name[int32(n)]
	return ok
}
