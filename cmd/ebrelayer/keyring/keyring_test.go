package keyring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKeyRing(t *testing.T) {

	mnemonic := "reject climb decline mule tell taste swing split pool stumble mask job offer exhaust bulk approve crawl alpha burst lion ribbon screen return have"
	moniker := "user"
	password := "12345678"

	keyRing := NewKeyRing(mnemonic, moniker, password)
	require.NotEqual(t, keyRing, nil, "keyring is nil")
}

func TestSign(t *testing.T) {

	mnemonic := "reject climb decline mule tell taste swing split pool stumble mask job offer exhaust bulk approve crawl alpha burst lion ribbon screen return have"
	moniker := "user"
	password := "12345678"

	keyRing := NewKeyRing(mnemonic, moniker, password)
	require.NotEqual(t, keyRing, nil, "keyring is nil")

	msg := []byte("Hello")

	_, _, err := keyRing.Sign(msg)
	require.NotEqual(t, err, nil, "failed to sign")
}
