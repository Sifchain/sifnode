import json
from typing import List, Optional


class Ganache:
    @staticmethod
    def start_ganache_cli(env, executable: Optional[str] = None, mnemonic: Optional[List[str]] = None,
        db: Optional[str] = None, port: Optional[int] = None, host: Optional[str] = None,
        network_id: Optional[int] = None, gas_price: Optional[int] = None, gas_limit: Optional[int] = None,
        default_balance_ether: Optional[int] = None, block_time: Optional[int] = None,
        account_keys_path: Optional[str] = None, log_file: Optional[str] = None
    ):
        args = \
            ([executable] if executable is not None else ["ganache-cli"]) + \
            (["--mnemonic", " ".join(mnemonic)] if mnemonic else []) + \
            (["--db", db] if db else []) + \
            (["--port", str(port)] if port is not None else []) + \
            (["--host", host] if host else []) + \
            (["--networkId", str(network_id)] if network_id is not None else []) + \
            (["--gasPrice", str(gas_price)] if gas_price is not None else []) + \
            (["--gasLimit", str(gas_limit)] if gas_limit is not None else []) + \
            (["--defaultBalanceEther", str(default_balance_ether)] if default_balance_ether is not None else []) + \
            (["--blockTime", str(block_time)] if block_time is not None else []) + \
            (["--account_keys_path", account_keys_path] if account_keys_path is not None else [])
        return env.popen(args, log_file=log_file)


class GanacheAbiProvider:
    def __init__(self, cmd, artifacts_dir, ethereum_network_id, deployed_smart_contract_address_overrides):
        self.cmd = cmd
        self.artifacts_dir = artifacts_dir
        self.ethereum_default_network_id = ethereum_network_id
        self.deployed_smart_contract_address_overrides = deployed_smart_contract_address_overrides

    def get_descriptor(self, sc_name):
        path = self.cmd.project.project_dir(self.artifacts_dir, "{}.json".format(sc_name))
        tmp = json.loads(self.cmd.read_text_file(path))
        abi = tmp["abi"]
        bytecode = tmp["bytecode"]
        deployed_address = None
        if (self.deployed_smart_contract_address_overrides is not None) and (sc_name in self.deployed_smart_contract_address_overrides):
            deployed_address = self.deployed_smart_contract_address_overrides[sc_name]
        else:
            if ("networks" in tmp) and (self.ethereum_default_network_id is not None):
                str_network_id = str(self.ethereum_default_network_id)
                if str_network_id in tmp["networks"]:
                    deployed_address = tmp["networks"][str_network_id]["address"]
        return abi, bytecode, deployed_address
