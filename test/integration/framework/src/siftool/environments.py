from siftool.common import *
from siftool.sifchain import ROWAN, STAKE
from siftool import sifchain, command, cosmos


# Environment for load test test_many_pools_and_liquidity_providers and for testing min commission/max voting power
# Just sifnode, no ethereum
# Multi-node support
class SifnodedEnvironment:
    def __init__(self, cmd: command.Command):
        self.cmd = cmd
        self.chain_id = "localnet"
        self.number_of_nodes = 1
        self.node_external_ip_address = LOCALHOST
        self.sifnoded_home_root = None
        self.log_level = None
        self.staking_denom = ROWAN
        self.faucet = None
        self.faucet_balance = {ROWAN: 10**30, STAKE: 10**30}
        # The stake of every validator must be greater than the balance of its administrator account..
        # At any time, 2/3 of validators (by stake) have to be online, otherwise there is no consensus and no new blocks.
        # We use the same stake for all validators.
        self.default_staking_amount = 92 * 10**21
        self.default_validator_balance = {ROWAN: 10**25}
        self.extra_genesis_balances = {}
        self.node_info = None
        self.sifnoded = None
        self.running_processes = None
        self.open_log_files = None
        self.default_commission_rate = 0.10
        self.default_commission_max_rate = 0.20
        self.default_commission_max_change_rate = 0.01
        self.default_min_self_delegation = 1000000
        self.default_initial_validator_mnemonic = None
        self.default_binary = None

    def init(self):
        assert self.default_commission_max_change_rate <= self.default_commission_max_rate, \
            "Commission max_change_rate cannot be more than commission_max_rate"

        self.node_info = []
        self.sifnoded = []
        self.sifnoded_home_root = self.sifnoded_home_root or self.cmd.mktempdir()

        for index in range(self.number_of_nodes):
            sifnoded, node_info = self._create_validator_home_and_account(index)
            self.sifnoded.append(sifnoded)
            self.node_info.append(node_info)

        sifnoded0 = self.sifnoded[0]
        self.faucet = sifnoded0.create_addr("faucet")

        for index in range(self.number_of_nodes):
            admin_addr = self.node_info[index]["admin_addr"]
            admin_bech = self.sifnoded[index].get_val_address(admin_addr)
            sifnoded0.add_genesis_account(admin_addr, cosmos.balance_add(self.default_validator_balance,
                {self.staking_denom: self.default_staking_amount}))
            sifnoded0.add_genesis_validators(admin_bech)
        admin0_addr = self.node_info[0]["admin_addr"]
        admin0_name = self.node_info[0]["admin_name"]
        sifnoded0.add_genesis_clp_admin(admin0_addr)
        sifnoded0.set_genesis_oracle_admin(admin0_name)
        sifnoded0.set_genesis_whitelister_admin(admin0_name)

        # Modify genesis.json of node0
        genesis = sifnoded0.load_genesis_json()
        all_genesis_balances = cosmos.balance_sum_by_address({self.faucet: self.faucet_balance}, self.extra_genesis_balances)
        if all_genesis_balances:
            sifnoded0.add_accounts_to_existing_genesis(genesis, all_genesis_balances)

        app_state = genesis["app_state"]
        app_state["gov"]["voting_params"] = {"voting_period": "60s"}
        # app_state["gov"]["deposit_params"]["min_deposit"] = [{"denom": ROWAN, "amount": "10000000"}]
        app_state["gov"]["deposit_params"]["min_deposit"] = [{"denom": STAKE, "amount": "10000000"}]
        app_state["crisis"]["constant_fee"] = {"denom": ROWAN, "amount": "1000"}
        app_state["staking"]["params"]["bond_denom"] = self.staking_denom
        app_state["mint"]["params"]["mint_denom"] = ROWAN
        sifnoded0.save_genesis_json(genesis)

        sifnoded0.gentx(admin0_name, {self.staking_denom: self.default_staking_amount},
            commission_rate=self.default_commission_rate, commission_max_rate=self.default_commission_max_rate,
            commission_max_change_rate=self.default_commission_max_change_rate)
        sifnoded0.collect_gentx()
        sifnoded0.validate_genesis()

        for index in range(self.number_of_nodes):
            self.update_configuration_files(index)

    def start(self):
        sifnoded0 = self.sifnoded[0]

        self.running_processes = []
        self.open_log_files = []
        for index in range(len(self.sifnoded)):
            log_file, process = self._sifnoded_start(index)
            self.running_processes.append(process)
            self.open_log_files.append(log_file)

        # Wait for some time so that nodes are fully booted
        sifnoded0.wait_for_last_transaction_to_be_mined()

        # Create a validator for all non-0 nodes. Node 0 needs to be up, but node i may or may not be up.
        # We're using the same stake for all non-0 nodes as for 0 node.
        for index in [i for i in range(self.number_of_nodes) if i != 0]:
            node_info = self.node_info[index]
            self._broadcast_create_validator_msg(node_info, self.default_staking_amount, self.default_commission_rate,
                self.default_commission_max_rate, self.default_commission_max_change_rate,
                self.default_min_self_delegation)

        # Sometimes the last validator is a bit slow in seeing itself added
        sifnoded0.wait_for_last_transaction_to_be_mined()
        assert all(len(self.sifnoded[i].query_staking_validators()) == self.number_of_nodes
            for i in range(self.number_of_nodes))

        # Do a dummy transfer of 1 rowan unit to check if transactions work
        self.fund(self.sifnoded[0].create_addr(), {ROWAN: 10 ** sifchain.ROWAN_DECIMALS})

    def _sifnoded_start(self, index: int):
        sifnoded = self.sifnoded[index]
        node_info = self.node_info[index]
        ports = node_info["ports"]
        log_file_path = os.path.join(sifnoded.home, "sifnoded.log")
        log_file = open(log_file_path, "w")
        self.open_log_files.append(log_file)
        process = sifnoded.sifnoded_start(log_file=log_file, log_level="debug", trace=True,
            tcp_url="tcp://{}:{}".format(ANY_ADDR, ports["rpc"]), p2p_laddr="{}:{}".format(ANY_ADDR, ports["p2p"]),
            grpc_address="{}:{}".format(ANY_ADDR, ports["grpc"]),
            grpc_web_address="{}:{}".format(ANY_ADDR, ports["grpc_web"]),
            address="tcp://{}:{}".format(ANY_ADDR, ports["address"]))
        sifnoded._wait_up()
        return log_file, process

    def fund(self, address: cosmos.Address, amounts: cosmos.Balance):
        return self.sifnoded[0].send_and_check(self.faucet, address, amounts)

    def add_validator(self, moniker: Optional[str] = None, staking_amount: Optional[int] = None,
        extra_funds: cosmos.Balance = None, commission_rate: Optional[float] = None,
        commission_max_rate: Optional[float] = None, commission_max_change_rate: Optional[float] = None,
        min_self_delegation: Optional[int] = None
    ) -> int:
        next_index = len(self.sifnoded)
        sifnoded, node_info = self._create_validator_home_and_account(next_index, moniker=moniker)
        self.node_info.append(node_info)
        self.sifnoded.append(sifnoded)
        self.update_configuration_files(next_index)

        admin_addr = node_info["admin_addr"]
        staking_amount = staking_amount if staking_amount is not None else self.default_staking_amount
        commission_rate = commission_rate if commission_rate is not None else self.default_commission_rate
        commission_max_rate = commission_max_rate if commission_max_rate is not None else self.default_commission_max_rate
        commission_max_change_rate = commission_max_change_rate if commission_max_change_rate is not None else self.default_commission_max_change_rate
        min_self_delegation = min_self_delegation if min_self_delegation is not None else self.default_min_self_delegation
        extra_funds = extra_funds if extra_funds is not None else self.default_validator_balance

        assert commission_max_change_rate <= commission_max_rate, \
            "Commission max_change_rate cannot be more than commission_max_rate"

        assert cosmos.balance_exceeds(extra_funds, {ROWAN: sifchain.sif_tx_fee_in_rowan}), \
            "Validator needs at least one sif_tx_fee_in_rowan to fund the transaction"
        staking_balance = {self.staking_denom: staking_amount}
        self.fund(admin_addr, cosmos.balance_add(extra_funds, staking_balance))

        # Start the newly added validator then broadcast the message "create validator" message.
        # In a real world scenario perhaps we would need to wait for the new validator to catch up before we add it?
        self.node_info.append(node_info)
        self.sifnoded.append(sifnoded)
        self._sifnoded_start(next_index)

        self._broadcast_create_validator_msg(node_info, staking_amount, commission_rate, commission_max_rate,
            commission_max_change_rate, min_self_delegation)
        return next_index

    def _broadcast_create_validator_msg(self, node_info: JsonDict, staking_amount: int, commission_rate: float,
        commission_max_rate: float, commission_max_change_rate, min_self_delegation: int
    ):
        stake = {self.staking_denom: staking_amount}
        admin_addr = node_info["admin_addr"]
        pubkey = node_info["pubkey"]
        moniker = node_info["moniker"]

        # Send "create validator" transaction. For this we need to use sifnoded with new validator's keystore, but with
        # "--node" pointing to existing (running) validator. We also check that the sender has enough balance for
        # staking and transaction itself.
        sifnoded_tmp = sifchain.Sifnoded(self.cmd, home=node_info["home"], chain_id=self.chain_id,
            node=self.node_info[0]["external_address"])

        validators_before = {v["description"]["moniker"] for v in sifnoded_tmp.query_staking_validators()}
        assert moniker not in validators_before

        admin_balance = sifnoded_tmp.get_balance(admin_addr)
        assert cosmos.balance_exceeds(admin_balance, {ROWAN: sifchain.sif_tx_fee_in_rowan}), \
            "Validator needs at least one sif_tx_fee_in_rowan to fund the transaction"
        assert cosmos.balance_exceeds(admin_balance, stake), \
            "Validator needs at least {} for staking".format(cosmos.balance_format(stake))

        res = sifnoded_tmp.staking_create_validator(stake, pubkey, moniker, commission_rate, commission_max_rate,
            commission_max_change_rate, min_self_delegation, admin_addr, broadcast_mode="block")
        sifchain.check_raw_log(res)

        validators_after = {v["description"]["moniker"] for v in sifnoded_tmp.query_staking_validators()}
        assert len(validators_after) == len(validators_before) + 1
        assert moniker in validators_after

    # Adjust configuration files for i != 0node.
    def update_configuration_files(self, index):
        sifnoded_i = self.sifnoded[index]
        node_info = self.node_info[index]
        # According to gzukel, nodes need just one peer to make sync work.
        # Star topology also makes it simpler to add additional nodes.
        peers = [sifchain.format_peer_address(node_info["node_id"], LOCALHOST, node_info["ports"]["p2p"])
            for node_info in [self.node_info[0]]]
        if index != 0:
            genesis = self.sifnoded[0].load_genesis_json()
            sifnoded_i.save_genesis_json(genesis)  # Copy genesis from validator 0 to all other
        app_toml = sifnoded_i.load_app_toml()
        config_toml = sifnoded_i.load_config_toml()
        app_toml["minimum-gas-prices"] = sif_format_amount(0.5, ROWAN)
        app_toml['api']['enable'] = True
        app_toml["api"]["address"] = sifchain.format_node_url(ANY_ADDR, node_info["ports"]["api"])
        config_toml["log_level"] = self.log_level  # TODO Probably redundant
        config_toml['p2p']["external_address"] = "{}:{}".format(self.node_external_ip_address, node_info["ports"]["p2p"])
        if index != 0:
            config_toml["p2p"]["persistent_peers"] = ",".join(peers)
        config_toml['p2p']['max_num_inbound_peers'] = 50
        config_toml['p2p']['max_num_outbound_peers'] = 50
        config_toml['p2p']['allow_duplicate_ip'] = True
        config_toml["rpc"]["pprof_laddr"] = "{}:{}".format(LOCALHOST, node_info["ports"]["pprof"])
        config_toml['moniker'] = node_info["moniker"]
        sifnoded_i.save_app_toml(app_toml)
        sifnoded_i.save_config_toml(config_toml)

    def _create_validator_home_and_account(self, next_index: int, moniker: Optional[str] = None):
        ports = self.ports_for_node(next_index)
        moniker = moniker or "sifnoded-{}".format(next_index)
        home = os.path.join(self.sifnoded_home_root, moniker)
        sifnoded_i = sifchain.Sifnoded(self.cmd, node=sifchain.format_node_url(ANY_ADDR, ports["rpc"]),
            home=home, chain_id=self.chain_id, binary=self.default_binary)
        admin_name = "admin-{}".format(next_index)
        mnemonic = None  # TODO
        admin_addr = sifnoded_i.create_addr(admin_name, mnemonic=mnemonic)
        sifnoded_i.init(moniker)
        node_id = sifnoded_i.tendermint_show_node_id()  # Taken from ${sifnoded_home}/config/node_key.json
        pubkey = sifnoded_i.tendermint_show_validator()  # Taken from ${sifnoded_home}/config/priv_validator_key.json
        node_info = {
            "moniker": moniker,
            "home": home,
            "node_id": node_id,
            "pubkey": pubkey,
            "admin_name": admin_name,
            "admin_addr": admin_addr,
            "ports": ports,
            "external_address": sifchain.format_node_url(self.node_external_ip_address, ports["rpc"])  # For --node
        }
        return sifnoded_i, node_info

    def ports_for_node(self, i: int) -> JsonDict:
        assert i < 10, "Change port configuration for 10 or more nodes"
        return {
            "p2p": 10276 + i,
            "grpc": 10909 + i,
            "grpc_web": 10919 + i,
            "address": 10276 + i,
            "rpc": 10286 + i,
            "api": 10131 + i,
            "pprof": 10606 + i,
        }
