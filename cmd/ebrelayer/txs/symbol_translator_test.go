package txs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	sifchainDenomFeedface = "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE"
	ethereumSymbolFeeface = "Face"
)

func TestNewSymbolTranslatorFromJsonBytes(t *testing.T) {
	_, err := NewSymbolTranslatorFromJsonBytes([]byte("foo"))
	assert.Error(t, err)

	q := ` {"ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE": "Face"} `
	x, err := NewSymbolTranslatorFromJsonBytes([]byte(q))
	assert.NoError(t, err)
	assert.NotNil(t, x)
	assert.Equal(t, x.SifchainToEthereum(sifchainDenomFeedface), ethereumSymbolFeeface)
	assert.Equal(t, x.EthereumToSifchain(ethereumSymbolFeeface), sifchainDenomFeedface)
	assert.Equal(t, x.SifchainToEthereum("verbatim"), "verbatim")
	assert.Equal(t, x.EthereumToSifchain("verbatim"), "verbatim")
}

func TestNewSymbolTranslator(t *testing.T) {
	s := NewSymbolTranslator()
	assert.Equal(t, s.SifchainToEthereum("something"), "something")
}
