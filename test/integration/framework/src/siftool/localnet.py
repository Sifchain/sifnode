import os
import json
from siftool.command import Command
from siftool.common import project_dir, siftool_logger


log = siftool_logger(__name__)


# This is called from run_env as a hook to run additional IBC chains defined in LOCALNET variable.
def run_localnet_hook():
    localnet_env_var = os.environ.get("LOCALNET")
    if not localnet_env_var:
        return

    localnet = Localnet()
    if not os.path.exists(localnet.node_module_dir):
        log.info("Installing localnet dependencies on first use in '{}'...".format(localnet.node_module_dir))
        localnet.install_deps()
    if not os.path.exists(localnet.bin_dir):
        log.info("Downloading localnet binaries on first use in '{}'...".format(localnet.bin_dir))
        localnet.download_binaries()

    if not os.path.exists(localnet.config_dir):
        log.info("Init all chains on first use in '{}'...".format(localnet.config_dir))
        localnet.init_all_chains()

    # rm -rf /tmp/localnet/config/cosmos/cosmoshub-testnet
    # mkdir -p /tmp/localnet/config/cosmos/cosmoshub-testnet
    # /tmp/localnet/bin/gaiad init cosmoshub-testnet --chain-id cosmoshub-testnet --home /tmp/localnet/config/cosmos/cosmoshub-testnet
    # /tmp/localnet/bin/gaiad keys add cosmos-validator --keyring-backend test --home /tmp/localnet/config/cosmos/cosmoshub-testnet
    # /tmp/localnet/bin/gaiad keys add cosmos-source --keyring-backend test --home /tmp/localnet/config/cosmos/cosmoshub-testnet
    # /tmp/localnet/bin/gaiad add-genesis-account cosmos-validator 10000000000000000000uphoton --keyring-backend test --home /tmp/localnet/config/cosmos/cosmoshub-testnet
    # /tmp/localnet/bin/gaiad add-genesis-account cosmos-source 10000000000000000000uphoton --keyring-backend test --home /tmp/localnet/config/cosmos/cosmoshub-testnet

    # rm -rf /tmp/localnet/config/sifchain/sifchain-testnet-1
    # mkdir -p /tmp/localnet/config/sifchain/sifchain-testnet-1
    # /tmp/localnet/bin/sifnoded init sifchain-testnet-1 --chain-id sifchain-testnet-1 --home /tmp/localnet/config/sifchain/sifchain-testnet-1
    # /tmp/localnet/bin/sifnoded keys add sifchain-validator --keyring-backend test --home /tmp/localnet/config/sifchain/sifchain-testnet-1
    # /tmp/localnet/bin/sifnoded keys add sifchain-source --keyring-backend test --home /tmp/localnet/config/sifchain/sifchain-testnet-1
    # /tmp/localnet/bin/sifnoded add-genesis-account sifchain-validator 10000000000000000000rowan --keyring-backend test --home /tmp/localnet/config/sifchain/sifchain-testnet-1
    # /tmp/localnet/bin/sifnoded add-genesis-account sifchain-source 10000000000000000000rowan --keyring-backend test --home /tmp/localnet/config/sifchain/sifchain-testnet-1

    # For each chain:
    # defaultGenesis = what was created in ${home}/config/genesis.json
    # remoteGenesis = curl (${node from config}/genesis).data e.g. https://rpc.testnet.cosmos.network:443/genesis
    # cleanedUpGenesis = cleanUpGenesisState(defaultGenesis, remoteGenesis)
    #
    # writeFile(genesis, "${home}/config/genesis.json")
    #
    # if sifchain: ${binPath}/${binary} set-gen-denom-whitelist ${home}/config/denoms.json --home ${home}
    #
    # ${binPath}/${binary} gentx ${validatorAccountName} ${amount}${denom} --chain-id ${chainId} --keyring-backend test --home ${home}
    # ${binPath}/${binary} collect-gentxs --home ${home}

    localnet.start_all_chains()  # Runs sifnoded and gaiad



def get_localnet_config(cmd):
    config = json.loads(cmd.read_text_file(cmd.project.project_dir("test/localnet/config/chains.json")))
    return config


def run(cmd, argv):
    log.debug(repr(argv))
    config = get_localnet_config(cmd)
    # Filter out items with "disabled": true
    config = {k: v for k, v in config.items() if not v.get("disabled", False)}
    tmpdir = cmd.mktempdir()
    log.debug(tmpdir)

    localnet = Localnet()
    localnet.init_all_chains()
    # localnet.start_all_chains()

    return


def download_ibc_binaries(cmd, chains_to_download=None, output_path=None):
    if not output_path:
        output_path = cmd.pwd()
    else:
        if not cmd.exists(output_path):
            cmd.mkdir(output_path)
    config = get_localnet_config(cmd)
    tmpdir = cmd.mktempdir()
    # We prefer to compile sifchain. Sentinel uses sourceUrl, but there is no Makefile.
    all_supported_chains = set(config.keys()).difference({"sifchain", "sentinel"})
    chains_to_download = chains_to_download or "all"
    if chains_to_download == "all":
        chains_to_download = all_supported_chains
    else:
        chains_to_download = ",".split(chains_to_download)
    try:
        tmp_gobin = os.path.join(tmpdir, "bin")
        cmd.mkdir(tmp_gobin)
        for chain_name in chains_to_download:
            if chain_name not in config:
                raise Exception("Chain {} not supported yet".format(chain_name))
            values = config[chain_name]
            binary = values["binary"]
            binary_url = values.get("binaryUrl")
            source_url = values.get("sourceUrl")
            binary_relative_path = values.get("binaryRelativePath")
            source_relative_path = values.get("sourceRelativePath")
            assert bool(source_url) ^ bool(binary_url)
            url = binary_url or source_url
            dlfile = os.path.join(tmpdir, "{}-download.tmp".format(chain_name))
            log.info("Downloading {} from '{}' to '{}'...".format(chain_name, url, dlfile))
            cmd.download_url(url, output_file=dlfile)
            extract_dir = os.path.join(tmpdir, chain_name)
            src_file = None
            cmd.mkdir(extract_dir)
            if url.endswith(".zip"):
                cmd.execst(["unzip", dlfile], cwd=extract_dir)
            elif url.endswith(".tar.gz"):
                cmd.execst(["tar", "xfz", dlfile], cwd=extract_dir)
            elif binary_url:
                # We have binaryUrl but it is not an archive => must be binary itself
                assert not source_url and not binary_relative_path
                src_file = dlfile
            if not src_file:
                if binary_url:
                    src_file = os.path.join(extract_dir, binary_relative_path if binary_relative_path else binary)
                if source_url:
                    src_dir = extract_dir if not source_relative_path else os.path.join(extract_dir, source_relative_path)
                    cmd.execst(["make", "install"], cwd=src_dir, env={"GOBIN": tmp_gobin})
                    src_file = os.path.join(tmp_gobin, binary)
            assert src_file
            dst_file = os.path.join(output_path, binary)
            cmd.copy_file(src_file, dst_file)
            cmd.chmod(dst_file, "+x")
    finally:
        cmd.rmf(tmpdir)


def fetch_genesis(base_url):
    pass


def init_chain(cmd):
    pass


class Localnet(Command):
    def __init__(self, script_dir=None, config_dir=None, bin_dir=None):
        self.script_dir = script_dir if script_dir else project_dir("test/localnet")
        self.config_dir = config_dir if config_dir else os.path.join("/tmp/localnet", "./config")
        self.bin_dir = bin_dir if bin_dir else os.path.join("/tmp/localnet", "./bin")
        self.node_module_dir = os.path.join(self.script_dir, "./node_modules")

    def install_deps(self):
        self.execst(["yarn"], cwd=self.script_dir, pipe=False)

    def download_binaries(self):
        self.execst(["yarn", "downloadBinaries"], cwd=self.script_dir, pipe=False)

    def init_all_chains(self):
        self.execst(["yarn", "initAllChains"], cwd=self.script_dir, pipe=False)

    def start_all_chains(self):
        self.execst(['yarn', 'startAllChains'], cmd=self.script_dir, pipe=False)

    def init_all_relayers(self):
        self.execst(['yarn', 'initAllRelayers'], cmd=self.script_dir, pipe=False)

    def start_all_relayers(self):
        self.execst(['yarn', 'startAllRelayers'], cmd=self.script_dir, pipe=False)

    def build_local_net(self):
        self.execst(['yarn', 'buildLocalNet'], cmd=self.script_dir, pipe=False)

    def load_local_net(self):
        self.execst(['yarn', 'loadLocalNet'], cmd=self.script_dir, pipe=False)

    def take_snapshot(self):
        self.execst(['yarn', 'takeSnapshot'], cmd=self.script_dir, pipe=False)

    def create_snapshot(self):
        self.execst(['yarn', 'createSnapshot'], cmd=self.script_dir, pipe=False)
        
    def test(self):
        self.execst(['yarn', 'test'], cmd=self.script_dir, pipe=False)
        