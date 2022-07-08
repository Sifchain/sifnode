package types

var WhitelistStorePrefix = []byte{0x01}

// We skip 0x02 because it was used by a prefix that has been deprecated
var TokenDenomPrefix = []byte{0x03}
