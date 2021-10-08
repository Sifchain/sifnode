package relayer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewKeybase test if we can get keybase from moniker, mnemonic and password
func TestNewKeybase(t *testing.T) {
	validatorMoniker := "akasha"
	mnemonic := "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
	password := ""

	base, info, err := NewKeybase(validatorMoniker, mnemonic, password)
	require.NotEqual(t, base, nil)
	require.NotEqual(t, info, nil)
	require.Equal(t, err, nil)
}
