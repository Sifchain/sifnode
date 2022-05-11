package types

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewNetworkIdentity get a new NetworkIdentity instance
func NewNetworkIdentity(networkDescriptor NetworkDescriptor) NetworkIdentity {
	return NetworkIdentity{
		NetworkDescriptor: networkDescriptor,
	}
}

// GetPrefix return storage prefix
func (n NetworkIdentity) GetPrefix(cdc codec.BinaryCodec) []byte {
	bytebuf := cdc.MustMarshal(&n)
	return append(WhiteListValidatorPrefix, bytebuf...)
}

// GetCrossChainFeePrefix return storage prefix
func (n NetworkIdentity) GetCrossChainFeePrefix(cdc codec.BinaryCodec) []byte {
	bytebuf := cdc.MustMarshal(&n)
	return append(CrossChainFeePrefix, bytebuf...)
}

// GetConsensusNeededPrefix return storage prefix
func (n NetworkIdentity) GetConsensusNeededPrefix(cdc codec.BinaryCodec) []byte {
	bytebuf := cdc.MustMarshal(&n)
	return append(ConsensusNeededPrefix, bytebuf...)
}

// GetFromPrefix return a NetworkIdentity from prefix
func GetFromPrefix(cdc codec.BinaryCodec, prefix []byte, key []byte) (NetworkIdentity, error) {
	if len(prefix) == 1 && len(key) >= 1 && prefix[0] == key[0] {
		var networkIdentity NetworkIdentity
		err := cdc.Unmarshal(key[1:], &networkIdentity)

		if err == nil {
			return networkIdentity, nil
		}
		return NetworkIdentity{}, err
	}

	return NetworkIdentity{}, errors.New("prefix is invalid")
}

// IsValid check if the network id is valid
func (n NetworkDescriptor) IsValid() bool {
	_, ok := NetworkDescriptor_name[int32(n)]
	return ok
}

// IsSifchain check if the network id is Sifchain
func (n NetworkDescriptor) IsSifchain() bool {
	return n == NetworkDescriptor_NETWORK_DESCRIPTOR_UNSPECIFIED
}
