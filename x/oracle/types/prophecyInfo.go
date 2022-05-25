package types

// GetSignaturePrefix return the signature key in keeper
func (info *ProphecyInfo) GetSignaturePrefix() []byte {
	return append(SignaturePrefix, info.ProphecyId...)
}
