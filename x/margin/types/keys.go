package types

var MTPPrefix = []byte{0x01}

func GetMTPKey(collateralAsset, custodyAsset, address string) []byte {
	return append(MTPPrefix,
		append([]byte(collateralAsset), append([]byte(custodyAsset), []byte(address)...)...)...,
	)
}
