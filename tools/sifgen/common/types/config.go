package types

type RPC struct {
	Laddr                     string   `toml:"laddr"`
	CorsAllowedOrigins        []string `toml:"cors_allowed_origins"`
	CorsAllowedMethods        []string `toml:"cors_allowed_methods"`
	CorsAllowedHeaders        []string `toml:"cors_allowed_headers"`
	GrpcLaddr                 string   `toml:"grpc_laddr"`
	GrpcMaxOpenConnections    int      `toml:"grpc_max_open_connections"`
	Unsafe                    bool     `toml:"unsafe"`
	MaxOpenConnections        int      `toml:"max_open_connections"`
	MaxSubscriptionClients    int      `toml:"max_subscription_clients"`
	MaxSubscriptionsPerClient int      `toml:"max_subscriptions_per_client"`
	TimeoutBroadcastTxCommit  string   `toml:"timeout_broadcast_tx_commit"`
	MaxBodyBytes              int      `toml:"max_body_bytes"`
	MaxHeaderBytes            int      `toml:"max_header_bytes"`
	TLSCertFile               string   `toml:"tls_cert_file"`
	TLSKeyFile                string   `toml:"tls_key_file"`
	PprofLaddr                string   `toml:"pprof_laddr"`
}

type P2P struct {
	Laddr                        string `toml:"laddr"`
	ExternalAddress              string `toml:"external_address"`
	Seeds                        string `toml:"seeds"`
	PersistentPeers              string `toml:"persistent_peers"`
	Upnp                         bool   `toml:"upnp"`
	AddrBookFile                 string `toml:"addr_book_file"`
	AddrBookStrict               bool   `toml:"addr_book_strict"`
	MaxNumInboundPeers           int    `toml:"max_num_inbound_peers"`
	MaxNumOutboundPeers          int    `toml:"max_num_outbound_peers"`
	UnconditionalPeerIds         string `toml:"unconditional_peer_ids"`
	PersistentPeersMaxDialPeriod string `toml:"persistent_peers_max_dial_period"`
	FlushThrottleTimeout         string `toml:"flush_throttle_timeout"`
	MaxPacketMsgPayloadSize      int    `toml:"max_packet_msg_payload_size"`
	SendRate                     int    `toml:"send_rate"`
	RecvRate                     int    `toml:"recv_rate"`
	Pex                          bool   `toml:"pex"`
	SeedMode                     bool   `toml:"seed_mode"`
	PrivatePeerIds               string `toml:"private_peer_ids"`
	AllowDuplicateIP             bool   `toml:"allow_duplicate_ip"`
	HandshakeTimeout             string `toml:"handshake_timeout"`
	DialTimeout                  string `toml:"dial_timeout"`
}

type Mempool struct {
	Recheck               bool   `toml:"recheck"`
	Broadcast             bool   `toml:"broadcast"`
	WalDir                string `toml:"wal_dir"`
	Size                  int    `toml:"size"`
	MaxTxsBytes           int    `toml:"max_txs_bytes"`
	CacheSize             int    `toml:"cache_size"`
	KeepInvalidTxsInCache bool   `toml:"keep-invalid-txs-in-cache"`
	MaxTxBytes            int    `toml:"max_tx_bytes"`
	MaxBatchBytes         int    `toml:"max_batch_bytes"`
}

type Statesync struct {
	Enable        bool   `toml:"enable"`
	RPCServers    string `toml:"rpc_servers"`
	TrustHeight   int    `toml:"trust_height"`
	TrustHash     string `toml:"trust_hash"`
	TrustPeriod   string `toml:"trust_period"`
	DiscoveryTime string `toml:"discovery_time"`
	TempDir       string `toml:"temp_dir"`
}

type Fastsync struct {
	Version string `toml:"version"`
}

type Consensus struct {
	WalFile                     string `toml:"wal_file"`
	TimeoutPropose              string `toml:"timeout_propose"`
	TimeoutProposeDelta         string `toml:"timeout_propose_delta"`
	TimeoutPrevote              string `toml:"timeout_prevote"`
	TimeoutPrevoteDelta         string `toml:"timeout_prevote_delta"`
	TimeoutPrecommit            string `toml:"timeout_precommit"`
	TimeoutPrecommitDelta       string `toml:"timeout_precommit_delta"`
	TimeoutCommit               string `toml:"timeout_commit"`
	DoubleSignCheckHeight       int    `toml:"double_sign_check_height"`
	SkipTimeoutCommit           bool   `toml:"skip_timeout_commit"`
	CreateEmptyBlocks           bool   `toml:"create_empty_blocks"`
	CreateEmptyBlocksInterval   string `toml:"create_empty_blocks_interval"`
	PeerGossipSleepDuration     string `toml:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration string `toml:"peer_query_maj23_sleep_duration"`
}

type TxIndex struct {
	Indexer string `toml:"indexer"`
}

type Instrumentation struct {
	Prometheus           bool   `toml:"prometheus"`
	PrometheusListenAddr string `toml:"prometheus_listen_addr"`
	MaxOpenConnections   int    `toml:"max_open_connections"`
	Namespace            string `toml:"namespace"`
}

type ConfigTOML struct {
	ProxyApp               string          `toml:"proxy_app"`
	Moniker                string          `toml:"moniker"`
	FastSync               bool            `toml:"fast_sync"`
	DbBackend              string          `toml:"db_backend"`
	DbDir                  string          `toml:"db_dir"`
	LogLevel               string          `toml:"log_level"`
	LogFormat              string          `toml:"log_format"`
	GenesisFile            string          `toml:"genesis_file"`
	PrivValidatorKeyFile   string          `toml:"priv_validator_key_file"`
	PrivValidatorStateFile string          `toml:"priv_validator_state_file"`
	PrivValidatorLaddr     string          `toml:"priv_validator_laddr"`
	NodeKeyFile            string          `toml:"node_key_file"`
	Abci                   string          `toml:"abci"`
	FilterPeers            bool            `toml:"filter_peers"`
	RPC                    RPC             `toml:"rpc"`
	P2P                    P2P             `toml:"p2p"`
	Mempool                Mempool         `toml:"mempool"`
	Statesync              Statesync       `toml:"statesync"`
	Fastsync               Fastsync        `toml:"fastsync"`
	Consensus              Consensus       `toml:"consensus"`
	TxIndex                TxIndex         `toml:"tx_index"`
	Instrumentation        Instrumentation `toml:"instrumentation"`
}
