import json
import re
from typing import Mapping, Any, Iterable
from siftool import eth, command
from siftool.common import *


def js_fmt(str, *params):
    esc_params = tuple(json.dumps(x) for x in params)
    return str.format(*esc_params)


# Documentation: https://geth.ethereum.org/docs/
# - Dev mode ("--dev") https://geth.ethereum.org/docs/getting-started/dev-mode
# - Private Network Tutorial: https://geth.ethereum.org/docs/getting-started/private-net
# - Private Networks: https://geth.ethereum.org/docs/interface/private-network
# - Running a standalone private Geth node for non-production purpose: https://medium.com/coinmonks/running-a-standalone-private-geth-node-for-non-production-purpose-d6e0ff226150

class Geth:
    def __init__(self, cmd: command.Command, datadir: Optional[str] = None):
        self.cmd = cmd
        self.program = "geth"
        self.datadir = datadir

    def geth_exec(self, geth_cmd_string, ipcpath):
        args = [self.program, "--exec", geth_cmd_string, ipcpath]
        return self.cmd.execst(args)

    class AttachEvalFunction:
        def __init__(self, geth, ipcpath):
            self.geth = geth
            self.ipcpath = ipcpath

        def __call__(self, js_expr, raw=False):
            args = [self.geth.program, "attach", "--exec", js_expr, self.ipcpath]
            res = stdout(self.geth.cmd.execst(args))
            return res if raw else json.loads(res)

        @property
        def coinbase_addr(self):
            js_expr = f"eth.coinbase"
            return self(js_expr)

        def create_account(self, password):
            js_expr = js_fmt("personal.newAccount({})", password)
            return self(js_expr)

        def unlock_account(self, addr, password):
            js_expr = js_fmt("personal.unlockAccount({}, {})", addr, password)
            # TODO Exception if unlock fails
            # Returns true if acount was unlocked successfully
            # Prints an error if not successful
            return self(js_expr)

        def get_balance(self, addr):
            js_expr = js_fmt("eth.getBalance({})", addr)
            return self(js_expr)

        # Amount is in wei
        # Returns txhash
        def send(self, from_addr, to_addr, amount):
            js_expr = js_fmt("eth.sendTransaction({{from: {}, to: {}, value: {}}})", from_addr, to_addr, amount)
            return self(js_expr)

    def attach_eval_fn(self, ipcpath):
        return Geth.AttachEvalFunction(self, ipcpath)

    # Creates a password-protected account in geth keyring for a given private key. This works deterministically,
    # meaning the account address/pubkey is the same for the same private key, and also the same that you would get
    # when creating address/pubkey in Hardhat.
    #
    # This uses "geth account import", the keys are stored in datadir/keys. The alternative is to use "geth console"
    # personal.createAccount(). This account will be visible in eth.accounts.
    #
    # Private key has is a hex string without "0x" prefix
    # Datadir cannot be the same datadir that a running geth uses
    # See "Creating an account by importing a private key": https://geth.ethereum.org/docs/interface/managing-your-accounts
    def create_account(self, private_key: eth.PrivateKey, password: Optional[str] = None) -> eth.Address:
        assert (not private_key.startswith("0x")) and (len(private_key) == 64)
        addr, key = eth.validate_address_and_private_key(None, private_key)
        passfile = self.cmd.mktempfile()
        keyfile = self.cmd.mktempfile()
        try:
            self.cmd.write_text_file(passfile, password or "")
            self.cmd.write_text_file(keyfile, private_key)
            args = [self.program, "account", "import", keyfile, "--password", passfile] + \
                (["--datadir", self.datadir] if self.datadir else [])
            res = self.cmd.execst(args)
            address = "0x" + re.compile("^Address: \\{(.*)\\}$").match(exactly_one(stdout_lines(res)))[1]
            assert addr.lower() == address
            return addr
        finally:
            self.cmd.rm(keyfile)
            self.cmd.rm(passfile)

    def create_genesis_config_clique(self, chain_id: int, signer_addresses: Iterable[eth.Address],
        alloc: Mapping[eth.Address, int], gas_limit: int = 8000000, difficulty: int = 1, block_mining_period: int = 5,
    ) -> Mapping[str, Any]:
        # See https://geth.ethereum.org/docs/interface/private-network
        extradata = "0x" + "00"*32 + ''.join([addr.lower()[2:] for addr in signer_addresses]) + "00"*65
        return {
            "config": {
                "chainId": chain_id,
                "homesteadBlock": 0,
                "eip150Block": 0,
                "eip155Block": 0,
                "eip158Block": 0,
                "byzantiumBlock": 0,
                "constantinopleBlock": 0,
                "petersburgBlock": 0,
                "istanbulBlock": 0,
                "berlinBlock": 0,
                "londonBlock": 0,
                "clique": {
                    "period": block_mining_period,
                    "epoch": 30000
                }
            },
            "difficulty": str(difficulty),
            "gasLimit": str(gas_limit),
            "extradata": extradata,
            "alloc": {k: {"balance": str(v)} for k, v in alloc.items()}
        }

    def init(self, ethereum_chain_id: int, signers: Iterable[eth.Address], gas_limit: Optional[int] = None,
        funds_alloc: Optional[Mapping[eth.Address, int]] = None, block_mining_period: Optional[int] = None
    ):
        funds_alloc = funds_alloc or {}
        kwargs = {}
        if block_mining_period is not None:
            kwargs["block_mining_period"] = block_mining_period
        if gas_limit is not None:
            kwargs["gas_limit"] = gas_limit
        tmp_genesis_file = self.cmd.mktempfile()
        try:
            genesis = self.create_genesis_config_clique(ethereum_chain_id, signers, funds_alloc, **kwargs)
            self.cmd.write_text_file(tmp_genesis_file, json.dumps(genesis))
            args = [self.program, "init", tmp_genesis_file] + (["--datadir", self.datadir] if self.datadir else [])
            self.cmd.execst(args)
        finally:
            self.cmd.rm(tmp_genesis_file)

    def buid_run_args(self, network_id: int, http_port: Optional[int] = None, ws_port: Optional[int] = None,
        dev: bool = False, mine: bool = False, unlock: Optional[Iterable[str]] = None, password: Optional[str] = None,
        allow_insecure_unlock: bool = False, rpc_allow_unprotected_txs: bool = False, gas_price: Optional[int] = None,
        gas_limit: Optional[int] = None, verbosity: Optional[int] = None
    ):
        all_apis = ["personal", "eth", "net", "web3"]
        args = self.build_geth_cmd_args(network_id=network_id, no_discover=True, dev=dev, mine=mine,
            miner_gas_price=gas_price, dev_gas_limit=gas_limit if dev else None,
            miner_gas_limit=gas_limit if not dev else None, unlock=unlock, password=password,
            allow_insecure_unlock=allow_insecure_unlock, datadir=self.datadir, ws=True, ws_addr=ANY_ADDR,
            ws_port=ws_port, ws_api=all_apis, http=True, http_addr=ANY_ADDR, http_port=http_port, http_api=all_apis,
            rpc_allow_unprotected_txs=rpc_allow_unprotected_txs, verbosity=verbosity)
        return command.buildcmd(args)

    def build_geth_cmd_args(self, command: Optional[str] = None, network_id: Optional[int] = None,
        datadir: Optional[str] = None, unlock: Optional[Iterable[str]] = None, password: Optional[str] = None,
        allow_insecure_unlock: bool = False, ipcpath: Optional[str] = None, ws: bool = False,
        ws_addr: Optional[str] = None, ws_port: Optional[int] = None, ws_api: Iterable[str] = None, http: bool = False,
        http_addr: Optional[str] = None, http_port: Optional[int] = None, http_api: Iterable[str] = None,
        rpc_allow_unprotected_txs: bool = False, dev: bool = False, dev_gas_limit: Optional[int] = None,
        dev_period: Optional[int] = None, rpcvhosts: Optional[str] = None, mine: bool = False,
        miner_gas_price: Optional[int] = None, miner_gas_limit: Optional[int] = None,
        miner_threads: Optional[int] = None, no_discover: bool = False, verbosity: Optional[str] = None
     ):
        args = [self.program] + \
            ([command] if command else []) + \
            (["--networkid", str(network_id)] if network_id else []) + \
            (["--datadir", datadir] if datadir else []) + \
            (["--unlock", ",".join(unlock)] if unlock else []) + \
            (["--password", password] if password else []) + \
            (["--allow-insecure-unlock"] if allow_insecure_unlock else []) + \
            (["--ipcpath", ipcpath] if ipcpath else []) + \
            (["--ws"] if ws else []) + \
            (["--ws.addr", ws_addr] if ws_addr else []) + \
            (["--ws.port", str(ws_port)] if ws_port is not None else []) + \
            (["--ws.api", ",".join(ws_api)] if ws_api else []) + \
            (["--http"] if http else []) + \
            (["--http.addr", http_addr] if http_addr is not None else []) + \
            (["--http.port", str(http_port)] if http_port is not None else []) + \
            (["--http.api", ",".join(http_api)] if http_api else []) + \
            (["--rpc.allow-unprotected-txs"] if rpc_allow_unprotected_txs else []) + \
            (["--dev"] if dev else []) + \
            (["--dev.gaslimit", str(dev_gas_limit)] if dev_gas_limit is not None else []) + \
            (["--dev.period", str(dev_period)] if dev_period is not None else []) + \
            (["--rpcvhosts", rpcvhosts] if rpcvhosts else []) + \
            (["--mine"] if mine else []) + \
            (["--miner.gasprice", str(miner_gas_price)] if miner_gas_price is not None else []) + \
            (["--miner.gaslimit", str(miner_gas_limit)] if miner_gas_limit is not None else []) + \
            (["--miner.threads", str(miner_threads)] if miner_threads is not None else []) + \
            (["--nodiscover"] if no_discover else []) + \
            (["--verbosity", str(verbosity)] if verbosity is not None else [])
        return args
