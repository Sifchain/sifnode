from siftool.common import *
from siftool.sifchain import ROWAN
from siftool import common, sifchain, command


# Environment for load test test_many_pools_and_liquidity_providers
# Just sifnode, no ethereum
# Multi-node support
class SifnodedEnvironment:
    def __init__(self, cmd: command.Command):
        self.cmd = cmd
        self.chain_id = None
        self.number_of_nodes = None
        self.node_external_ip_address = None
        self.sifnoded_home_root = None
        self.validator0_mnemonic = None
        self.log_level = None
        self.validator_account_balance = None
        self.genesis_balances = None
        self.node_info = None
        self.sifnoded = None
        self.running_processes = None
        self.open_log_files = None

    def ports_for_node(self, i: int) -> JsonDict:
        return {
            "p2p": 10276 + i,
            "grpc": 10909 + i,
            "grpc_web": 10919 + i,
            "address": 10276 + i,
            "rpc": 10286 + i,
            "api": 10131 + i,
            "pprof": 10606 + i,
        }

    def init(self):
        self.sifnoded = []
        self.node_info = []
        for i in range(self.number_of_nodes):
            ports = self.ports_for_node(i)
            home = os.path.join(self.sifnoded_home_root, "sifnoded-{}".format(i))
            sifnoded_i = sifchain.Sifnoded(self.cmd, node=sifchain.format_node_url(ANY_ADDR, ports["rpc"]),
                home=home, chain_id=self.chain_id)
            moniker = "sifnode-{}".format(i)
            acct_name = "sif-{}".format(i)
            acct_addr = sifnoded_i.create_addr(acct_name, mnemonic=self.validator0_mnemonic if i == 0 else None)
            sifnoded_i.init(moniker)
            node_id = sifnoded_i.tendermint_show_node_id()  # Taken from ${sifnoded_home}/config/node_key.json
            pubkey = sifnoded_i.tendermint_show_validator()  # Taken from ${sifnoded_home}/config/priv_validator_key.json
            node_info = {
                "moniker": moniker,
                "home": home,
                "node_id": node_id,
                "pubkey": pubkey,
                "acct_name": acct_name,
                "acct_addr": acct_addr,
                "ports": ports,
                "external_address": sifchain.format_node_url(self.node_external_ip_address, ports["rpc"])  # For --node
            }
            self.sifnoded.append(sifnoded_i)
            self.node_info.append(node_info)

        sifnoded0 = self.sifnoded[0]

        for i in range(self.number_of_nodes):
            acct_addr = self.node_info[i]["acct_addr"]
            acct_bech = self.sifnoded[i].get_val_address(acct_addr)
            sifnoded0.add_genesis_validators(acct_bech)
            if self.validator_account_balance:
                sifnoded0.add_genesis_account(acct_addr, self.validator_account_balance)
        admin0_addr = self.node_info[0]["acct_addr"]
        admin0_name = self.node_info[0]["acct_name"]
        sifnoded0.add_genesis_clp_admin(admin0_addr)
        sifnoded0.set_genesis_oracle_admin(admin0_name)
        sifnoded0.set_genesis_whitelister_admin(admin0_name)

        genesis = sifnoded0.load_genesis_json()
        if self.genesis_balances:
            sifnoded0.add_accounts_to_existing_genesis(genesis, self.genesis_balances)

        app_state = genesis["app_state"]
        app_state["gov"]["voting_params"] = {"voting_period": "120s"}
        app_state["gov"]["deposit_params"]["min_deposit"] = [{"denom": ROWAN, "amount": "10000000"}]
        app_state["crisis"]["constant_fee"] = {"denom": ROWAN, "amount": "1000"}
        app_state["staking"]["params"]["bond_denom"] = ROWAN
        app_state["mint"]["params"]["mint_denom"] = ROWAN
        sifnoded0.save_genesis_json(genesis)

        sifnoded0.gentx(admin0_name, {ROWAN: 10**24})
        sifnoded0.collect_gentx()
        sifnoded0.validate_genesis()

        # According to gzukel, nodes need just one peer to make sync work.
        peers = [sifchain.format_peer_address(node_info["node_id"], LOCALHOST, node_info["ports"]["p2p"])
            for node_info in [self.node_info[0]]]
        genesis = sifnoded0.load_genesis_json()
        for i in range(self.number_of_nodes):
            sifnoded_i = self.sifnoded[i]
            if i != 0:
                sifnoded_i.save_genesis_json(genesis)  # Copy genesis from validator 0 to all other
            info = self.node_info[i]
            app_toml = sifnoded_i.load_app_toml()
            app_toml["minimum-gas-prices"] = sif_format_amount(0.5, ROWAN)
            app_toml['api']['enable'] = True
            app_toml["api"]["address"] = sifchain.format_node_url(ANY_ADDR, info["ports"]["api"])
            sifnoded_i.save_app_toml(app_toml)
            config_toml = sifnoded_i.load_config_toml()
            config_toml["log_level"] = self.log_level  # TODO Probably redundant
            config_toml['p2p']["external_address"] = "{}:{}".format(self.node_external_ip_address, info["ports"]["p2p"])
            if i != 0:
                config_toml["p2p"]["persistent_peers"] = ",".join(peers)
            config_toml['p2p']['max_num_inbound_peers'] = 50
            config_toml['p2p']['max_num_outbound_peers'] = 50
            config_toml['p2p']['allow_duplicate_ip'] = True
            config_toml["rpc"]["pprof_laddr"] = "{}:{}".format(LOCALHOST, info["ports"]["pprof"])
            config_toml['moniker'] = info["moniker"]
            sifnoded_i.save_config_toml(config_toml)

        # Start processes
        self.running_processes = []
        self.open_log_files = []
        for i, sifnoded_i in enumerate(self.sifnoded):
            node_info = self.node_info[i]
            ports = node_info["ports"]
            log_file_path = os.path.join(sifnoded_i.home, "sifnoded.log")
            log_file = open(log_file_path, "w")
            self.open_log_files.append(log_file)
            process = sifnoded_i.sifnoded_start(log_file=log_file, log_level="debug", trace=True,
                tcp_url="tcp://{}:{}".format(ANY_ADDR, ports["rpc"]), p2p_laddr="{}:{}".format(ANY_ADDR, ports["p2p"]),
                grpc_address="{}:{}".format(ANY_ADDR, ports["grpc"]),
                grpc_web_address="{}:{}".format(ANY_ADDR, ports["grpc_web"]),
                address="tcp://{}:{}".format(ANY_ADDR, ports["address"])
            )
            sifnoded_i._wait_up()
            self.running_processes.append(process)

        # Wait for some time so that nodes are fully booted
        sifnoded0.wait_for_last_transaction_to_be_mined()

        # Create a validator for all non-0 nodes. Node 0 needs to be up, but node i may or may not be up.
        for i in [x for x in range(self.number_of_nodes) if x != 0]:
            node_info = self.node_info[i]
            # This needs to have the private key ("home") of i-th validator but "node" of the 0-th.
            # TODO We need to use "rpc" for --node, not p2p / external_address!
            sifnoded_tmp = sifchain.Sifnoded(self.cmd, home=node_info["home"], chain_id=self.chain_id,
                node=self.node_info[0]["external_address"])
            sifnoded_tmp.staking_create_validator((10 ** 24, ROWAN), node_info["pubkey"], node_info["moniker"],
                0.10, 0.20, 0.01, 1000000, node_info["acct_addr"])

        sifnoded0.wait_for_last_transaction_to_be_mined()
