package types

var MTPPrefix = []byte{0x01}

func GetMTPKey(collateralAsset, custodyAsset, address string, position Position) []byte {
	return append(MTPPrefix,
		append([]byte(collateralAsset), append([]byte(custodyAsset), append([]byte(address), []byte(position.String())...)...)...)...,
	)
}
