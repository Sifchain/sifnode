package txs

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/utilities"
	"github.com/vishalkuo/bimap"
	"io/ioutil"
)

// SymbolTranslator translates between Sifchain denoms and Ethereum symbols
type SymbolTranslator struct {
	symbolTable *bimap.BiMap
}

func (s *SymbolTranslator) SifchainToEthereum(denom string) string {
	return utilities.GetWithDefault(s.symbolTable, denom, denom).(string)
}

func (s *SymbolTranslator) EthereumToSifchain(symbol string) string {
	return utilities.GetInverseWithDefault(s.symbolTable, symbol, symbol).(string)
}

func NewSymbolTranslator() *SymbolTranslator {
	return &SymbolTranslator{symbolTable: bimap.NewBiMap()}
}

func NewSymbolTranslatorFromJsonFile(filename string) (*SymbolTranslator, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewSymbolTranslatorFromJsonBytes(contents)
}

func NewSymbolTranslatorFromJsonBytes(jsonContents []byte) (*SymbolTranslator, error) {
	var symbolMap map[string]interface{}
	err := json.Unmarshal(jsonContents, &symbolMap)
	if err != nil {
		return nil, err
	}
	symbolBiMap := bimap.NewBiMap()
	for k, v := range symbolMap {
		symbolBiMap.Insert(k, v)
	}
	symbolBiMap.MakeImmutable()
	return &SymbolTranslator{symbolTable: symbolBiMap}, nil
}
