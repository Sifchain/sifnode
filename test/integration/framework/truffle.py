class Ganache:
    @staticmethod
    def start_ganache_cli(env, mnemonic=None, db=None, port=None, host=None, network_id=None, gas_price=None,
        gas_limit=None, default_balance_ether=None, block_time=None, account_keys_path=None, log_file=None
    ):
        args = ["ganache-cli"] + \
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
