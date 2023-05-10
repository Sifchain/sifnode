import json
import web3
from siftool import eth, command
from siftool.command import Command, ExecResult, buildcmd
from siftool.common import *


class Hardhat:
    def __init__(self, cmd: Command):
        assert on_peggy2_branch
        self.cmd = cmd
        self.project = cmd.project

    def build_start_args(self, hostname=None, port=None, fork=None, fork_block_number=None):
        # TODO We need to manaege smart-contracts/hardhat.config.ts + it also reads smart-contracts/.env via dotenv
        # TODO Handle failures, e.g. if the process is already running we get exit value 1 and
        # "Error: listen EADDRINUSE: address already in use 127.0.0.1:8545"
        args = [os.path.join("node_modules", ".bin", "hardhat"), "node"] + \
            (["--hostname", hostname] if hostname else []) + \
            (["--port", str(port)] if port is not None else []) + \
            (["--fork", fork] if fork else []) + \
            (["--fork-block-number", str(fork_block_number)] if fork_block_number is not None else [])
        return buildcmd(args, cwd=self.project.smart_contracts_dir)

    def compile_smart_contracts(self):
        # Creates:
        # smart-contracts/artifacts
        # smart-contracts/build
        # smart-contracts/cache
        self.project.npx(["hardhat", "compile"], cwd=project_dir("smart-contracts"), pipe=False)

    def script_runner(self, url: str = None, network: Optional[str] = None,
        ethereum_private_key: Optional[eth.PrivateKey] = None, accounts: Optional[Sequence[eth.PrivateKey]] = None
    ) -> 'ScriptRunner':
        return ScriptRunner(self, url=url, network=network, ethereum_private_key=ethereum_private_key, accounts=accounts)


# A wrapper around "npx hardhat run" for running TypeScript scripts that use hardhat.config.ts and need certain
# parameters/environment variables. This ensures we run all the scripts in a consistent way.
class ScriptRunner:
    def __init__(self, hardhat: Hardhat, url: str = None, network: Optional[str] = None,
        ethereum_private_key: Optional[eth.PrivateKey] = None, accounts: Optional[Sequence[eth.PrivateKey]] = None
    ):
        self.hardhat = hardhat
        self.url = url
        self.network = network
        self.ethereum_private_key = ethereum_private_key
        self.accounts = accounts

    # Values for 'network' parameter:
    # (see https://hardhat.org/getting-started/#connecting-a-wallet-or-dapp-to-hardhat-network):
    # - None:
    # - "localhost": connect to "http://127.0.0.1:8545" where "npx hardhat node" is running
    # - anything else: use the corresponding section from smart-contracts/hardhat.config.ts, element "networks"
    def deploy_smart_contracts(self) -> Mapping[str, eth.Address]:
        stdout_lines = self.run("deploy_contracts_dev.ts")
        tmp = json.loads(stdout_lines[-1])
        return {
            "BridgeBank": tmp["bridgeBank"],
            "BridgeRegistry": tmp["bridgeRegistry"],
            "CosmosBridge": tmp["cosmosBridge"],
            "Rowan": tmp["rowanContract"],
            "Blocklist": tmp["blocklist"],
        }

    # TODO This is called mostly from siftool and it's one line, convert to web3 call
    def update_validator_power(self, cosmos_bridge_addr, evm_validator_addresses, sifnode_witnesses):
        npx_env = {
            "COSMOSBRIDGE": cosmos_bridge_addr,
            "POWERS": ",".join(str(x["power"]) for x in sifnode_witnesses),
            "VALIDATORS": ",".join(evm_validator_addresses)
        }
        self.run("update_validator_power.ts", npx_env=npx_env)

    def run(self, script: str, npx_env: Optional[Mapping[str, str]] = None) -> ExecResult:
        args = ["hardhat", "run", "scripts/{}".format(script)] + \
            (["--network", self.network] if self.network else [])
        res = self.__npx_hardhat(args, npx_env=npx_env)
        # Skip first line "No need to generate any newer types". This only works if the smart contracts have already
        # been compiled, otherwise the output starts with 4 lines:
        #     Compiling 35 files with 0.5.16
        #     Generating typings for: 36 artifacts in dir: build for target: ethers-v5
        #     Successfully generated 65 typings!
        #     Compilation finished successfully
        # With devtool, the compilation is performed automatically before invoking main() if the script is invoked
        # via "npx hardhat run scripts/devenv.ts" instead of "npx ts-node scripts/devenv.ts", so normally this would
        # not happen.
        # TODO Suggested solution: pass a parameter to deploy_contracts.ts where it should write the output json file
        stdout_lines = stdout(res).splitlines()
        assert len(stdout_lines) > 0
        assert stdout_lines[0] == "No need to generate any newer typings."  # This is printed by hardhat, the rest is from the script
        return stdout_lines[1:]

    def test(self, test_files: Sequence[str], npx_env: Optional[Mapping[str, str]] = None):
        args = ["hardhat", "test" ] + list(test_files) + \
            (["--network", self.network] if self.network else [])
        return self.__npx_hardhat(args, npx_env=npx_env, pipe=False)

    def __npx_hardhat(self, args, npx_env: Optional[Mapping[str, str]] = None, pipe: bool = True):
        # If this fails with tsyringe complaining about missing "../../build" directory, do this:
        # rm -rf smart-contracts/artifacts.
        env = {}
        if self.url:
            env["NETWORK_URL"] = self.url
        if self.accounts:
            env["ETH_ACCOUNTS"] = ",".join(self.accounts)
        if self.ethereum_private_key:
            env["ETHEREUM_PRIVATE_KEY"] = self.ethereum_private_key
        if npx_env:
            env = dict_merge(env, npx_env)
        if not env:
            # Avoid passing empty environment, it crashes hardhat on some string split presumably because a variable
            # which it expects does not exist.
            env = None
        return self.hardhat.project.npx(args, cwd=self.hardhat.project.smart_contracts_dir, env=env, pipe=pipe)

class HardhatAbiProvider:
    def __init__(self, cmd: Command, abi_files_root: str, deployed_contract_addresses: Mapping[str, eth.Address]):
        self.cmd = cmd
        self.abi_files_root = abi_files_root
        self.deployed_contract_addresses = deployed_contract_addresses

    def get_descriptor(self, sc_name):
        relpath = {
            "BridgeBank": ["BridgeBank"],
            "BridgeToken": ["BridgeBank"],
            "CosmosBridge": [],
            "Rowan": ["BridgeBank"],
            "TrollToken": ["Mocks"],
            "FailHardToken": ["Mocks"],
            "UnicodeToken": ["Mocks"],
            "CommissionToken": ["Mocks"],
            "RandomTrollToken": ["Mocks"],
        }.get(sc_name, []) + [f"{sc_name}.sol", f"{sc_name}.json"]
        path = os.path.join(self.abi_files_root, *relpath)
        tmp = json.loads(self.cmd.read_text_file(path))
        abi = tmp["abi"]
        bytecode = tmp["bytecode"]
        deployed_address = self.deployed_contract_addresses.get(sc_name)
        return abi, bytecode, deployed_address


def default_accounts():
    # Hardhat doesn't provide a way to get the private keys of its default accounts, so just hardcode them for now.
    # TODO hardhat prints 20 accounts upon startup
    # Keep synced to smart-contracts/src/devenv/hardhatNode.ts:defaultHardhatAccounts
    # Format: [address, private_key]
    # Note: for compatibility with ganache, private keys should be stripped of "0x" prefix
    # (when you pass a private key to ebrelayer via ETHEREUM_PRIVATE_KEY, the key is treated as invalid)
    return [[web3.Web3.to_checksum_address(address), private_key] for address, private_key in [[
        "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
        "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    ], [
        "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
        "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
    ], [
        "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc",
        "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
    ], [
        "0x90f79bf6eb2c4f870365e785982e1f101e93b906",
        "7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    ], [
        "0x15d34aaf54267db7d7c367839aaf71a00a2c6a65",
        "47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
    ], [
        "0x9965507d1a55bcc2695c58ba16fb37d819b0a4dc",
        "8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
    ], [
        "0x976ea74026e726554db657fa54763abd0c3a0aa9",
        "92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
    ], [
        "0x14dc79964da2c08b23698b3d3cc7ca32193d9955",
        "4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
    ], [
        "0x23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f",
        "dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
    ], [
        "0xa0ee7a142d267c1f36714e4a8f75612f20a79720",
        "2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
    ]]]
