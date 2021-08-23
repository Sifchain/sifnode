package utils

import (
	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

const (
	TestMnemonic = "name chaos angry battle goat roast cause south wisdom creek trade regret fluid broccoli remind charge peasant photo option stock scene video baby clerk"
	TestPassword = "deepdarksecret"
)

func SafeTempDir(dir string, pattern string, t *testing.T) string {
	result, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		t.Errorf("failed to create temporary directory %v with pattern %v", dir, pattern)
	}
	return result
}

func TestCLI_AddKeyBackendFile(t *testing.T) {
	AddKeyToBackend(keyring.BackendFile, t)
}

func TestCLI_AddKeyBackendTest(t *testing.T) {
	AddKeyToBackend(keyring.BackendTest, t)
}

func TestCLI_AddKeyBackendMemory(t *testing.T) {
	AddKeyToBackend(keyring.BackendTest, t)
}

func AddKeyToBackend(backend string, t *testing.T) {
	c := NewCLI("unimportant", backend)
	tempdir := SafeTempDir("", "homedir", t)
	keyname := "mykey"
	key, err := c.AddKey(keyname, TestMnemonic, TestPassword, tempdir)
	require.NoError(t, err)
	var keyResult common.Keys
	unmarshallResult := yaml.Unmarshal([]byte(*key), &keyResult)
	if unmarshallResult != nil {
		t.Errorf("Failed to add key, error is: %v, input is: %v", unmarshallResult, *key)
	}
	if keyResult[0].Name != keyname {
		t.Errorf("Key was not added correctly")
	}
}
