package types

type Telemetry struct {
	ServiceName             string        `toml:"service-name"`
	Enabled                 bool          `toml:"enabled"`
	EnableHostname          bool          `toml:"enable-hostname"`
	EnableHostnameLabel     bool          `toml:"enable-hostname-label"`
	EnableServiceLabel      bool          `toml:"enable-service-label"`
	PrometheusRetentionTime int           `toml:"prometheus-retention-time"`
	GlobalLabels            []interface{} `toml:"global-labels"`
}

type API struct {
	Enable             bool   `toml:"enable"`
	Swagger            bool   `toml:"swagger"`
	Address            string `toml:"address"`
	MaxOpenConnections int    `toml:"max-open-connections"`
	RPCReadTimeout     int    `toml:"rpc-read-timeout"`
	RPCWriteTimeout    int    `toml:"rpc-write-timeout"`
	RPCMaxBodyBytes    int    `toml:"rpc-max-body-bytes"`
	EnabledUnsafeCors  bool   `toml:"enabled-unsafe-cors"`
}

type Grpc struct {
	Enable  bool   `toml:"enable"`
	Address string `toml:"address"`
}

type StateSync struct {
	SnapshotInterval   int `toml:"snapshot-interval"`
	SnapshotKeepRecent int `toml:"snapshot-keep-recent"`
}

type AppTOML struct {
	MinimumGasPrices  string        `toml:"minimum-gas-prices"`
	Pruning           string        `toml:"pruning"`
	PruningKeepRecent string        `toml:"pruning-keep-recent"`
	PruningKeepEvery  string        `toml:"pruning-keep-every"`
	PruningInterval   string        `toml:"pruning-interval"`
	HaltHeight        int           `toml:"halt-height"`
	HaltTime          int           `toml:"halt-time"`
	MinRetainBlocks   int           `toml:"min-retain-blocks"`
	InterBlockCache   bool          `toml:"inter-block-cache"`
	IndexEvents       []interface{} `toml:"index-events"`
	Telemetry         Telemetry     `toml:"telemetry"`
	API               API           `toml:"api"`
	Grpc              Grpc          `toml:"grpc"`
	StateSync         StateSync     `toml:"state-sync"`
}
