package types

type Config struct {
	ChainID        string `toml:"chain-id"`
	KeyringBackend string `toml:"keyring-backend"`
	Indent         bool   `toml:"indent"`
	TrustNode      bool   `toml:"trust-node"`
}
