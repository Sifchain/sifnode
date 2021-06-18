package types

type NodeConfig struct {
	FilterPeers            bool            `toml:"filter_peers"`
	FastSync               bool            `toml:"fast_sync"`
	ProxyApp               string          `toml:"proxy_app"`
	Moniker                string          `toml:"moniker"`
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
	ProfLaddr              string          `toml:"prof_laddr"`
	Fastsync               Fastsync        `toml:"fastsync"`
	P2P                    P2P             `toml:"p2p"`
	RPC                    RPC             `toml:"rpc"`
	Consensus              Consensus       `toml:"consensus"`
	Mempool                Mempool         `toml:"mempool"`
	Instrumentation        Instrumentation `toml:"instrumentation"`
	TxIndex                TxIndex         `toml:"tx_index"`
}

type RPC struct {
	Unsafe                    bool     `toml:"unsafe"`
	GrpcMaxOpenConnections    int      `toml:"grpc_max_open_connections"`
	MaxOpenConnections        int      `toml:"max_open_connections"`
	MaxSubscriptionClients    int      `toml:"max_subscription_clients"`
	MaxSubscriptionsPerClient int      `toml:"max_subscriptions_per_client"`
	MaxBodyBytes              int      `toml:"max_body_bytes"`
	MaxHeaderBytes            int      `toml:"max_header_bytes"`
	TimeoutBroadcastTxCommit  string   `toml:"timeout_broadcast_tx_commit"`
	Laddr                     string   `toml:"laddr"`
	GrpcLaddr                 string   `toml:"grpc_laddr"`
	TLSCertFile               string   `toml:"tls_cert_file"`
	TLSKeyFile                string   `toml:"tls_key_file"`
	CorsAllowedOrigins        []string `toml:"cors_allowed_origins"`
	CorsAllowedMethods        []string `toml:"cors_allowed_methods"`
	CorsAllowedHeaders        []string `toml:"cors_allowed_headers"`
}

type P2P struct {
	Upnp                         bool   `toml:"upnp"`
	AddrBookStrict               bool   `toml:"addr_book_strict"`
	Pex                          bool   `toml:"pex"`
	SeedMode                     bool   `toml:"seed_mode"`
	AllowDuplicateIP             bool   `toml:"allow_duplicate_ip"`
	MaxNumInboundPeers           int    `toml:"max_num_inbound_peers"`
	MaxNumOutboundPeers          int    `toml:"max_num_outbound_peers"`
	MaxPacketMsgPayloadSize      int    `toml:"max_packet_msg_payload_size"`
	SendRate                     int    `toml:"send_rate"`
	RecvRate                     int    `toml:"recv_rate"`
	Laddr                        string `toml:"laddr"`
	ExternalAddress              string `toml:"external_address"`
	Seeds                        string `toml:"seeds"`
	PersistentPeers              string `toml:"persistent_peers"`
	AddrBookFile                 string `toml:"addr_book_file"`
	UnconditionalPeerIds         string `toml:"unconditional_peer_ids"`
	PersistentPeersMaxDialPeriod string `toml:"persistent_peers_max_dial_period"`
	FlushThrottleTimeout         string `toml:"flush_throttle_timeout"`
	PrivatePeerIds               string `toml:"private_peer_ids"`
	HandshakeTimeout             string `toml:"handshake_timeout"`
	DialTimeout                  string `toml:"dial_timeout"`
}

type Mempool struct {
	Recheck     bool   `toml:"recheck"`
	Broadcast   bool   `toml:"broadcast"`
	WalDir      string `toml:"wal_dir"`
	Size        int    `toml:"size"`
	MaxTxsBytes int    `toml:"max_txs_bytes"`
	CacheSize   int    `toml:"cache_size"`
	MaxTxBytes  int    `toml:"max_tx_bytes"`
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
	SkipTimeoutCommit           bool   `toml:"skip_timeout_commit"`
	CreateEmptyBlocks           bool   `toml:"create_empty_blocks"`
	CreateEmptyBlocksInterval   string `toml:"create_empty_blocks_interval"`
	PeerGossipSleepDuration     string `toml:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration string `toml:"peer_query_maj23_sleep_duration"`
}

type TxIndex struct {
	IndexAllKeys bool   `toml:"index_all_keys"`
	Indexer      string `toml:"indexer"`
	IndexKeys    string `toml:"index_keys"`
}

type Instrumentation struct {
	Prometheus           bool   `toml:"prometheus"`
	PrometheusListenAddr string `toml:"prometheus_listen_addr"`
	MaxOpenConnections   int    `toml:"max_open_connections"`
	Namespace            string `toml:"namespace"`
}
