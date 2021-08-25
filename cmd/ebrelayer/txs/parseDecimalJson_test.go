package txs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDecimal(t *testing.T) {
	result := ParseDecimalFile("cosmosDecimalMap.json")
	expected := map[string]int{}
	expected["rowan"] = -8
	expected["cosmos_stake"] = -7

	require.Equal(t, result, expected)
}
