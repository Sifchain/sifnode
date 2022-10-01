package symbol_translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sifchainDenomFeedface = "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE"
	ethereumSymbolFeeface = "Face"
)

func TestNewSymbolTranslatorFromJsonBytes(t *testing.T) {
	_, err := NewSymbolTranslatorFromJSONBytes([]byte("foo"))
	assert.Error(t, err)

	q := ` {"ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE": "Face"} `
	x, err := NewSymbolTranslatorFromJSONBytes([]byte(q))
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
