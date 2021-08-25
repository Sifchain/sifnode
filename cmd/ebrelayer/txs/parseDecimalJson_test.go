package txs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDecimal(t *testing.T) {
	result := ParseDecimalFile("ethereumDecimalMap.json")
	expected := map[string]int{}
	expected["uatom"] = -12
	expected["uakt"] = -12
	expected["udvpn"] = -12

	require.Equal(t, result, expected)
}
