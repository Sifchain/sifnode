import os
import json
from siftool.command import Command, ExecResult
from siftool.common import *


log = siftool_logger(__name__)


def force_kill_processes(cmd):
    cmd.execst(["pkill", "node"], check_exit=False)
    cmd.execst(["pkill", "ebrelayer"], check_exit=False)
    cmd.execst(["pkill", "sifnoded"], check_exit=False)

def killall(processes):
    # TODO Order - ebrelayer, sifnoded, ganache
    for p in processes:
        if p is not None:
            p.kill()
            p.wait()


class Project:
    """Represents a checked out copy of a project in a particular directory."""

    def __init__(self, cmd: Command, base_dir: str):
        self.cmd = cmd
        self.base_dir = base_dir
        self.smart_contracts_dir = project_dir("smart-contracts")
        self.test_integration_dir = project_dir("test", "integration")
        self.siftool_dir = project_dir("test", "integration", "framework")
        self.go_path = os.environ.get("GOPATH")
        if self.go_path is None:
            # https://pkg.go.dev/cmd/go#hdr-GOPATH_and_Modules
            self.go_path = os.path.join(os.environ["HOME"], "go")
        self.go_bin_dir = os.environ.get("GOBIN")
        if self.go_bin_dir is None:
            self.go_bin_dir = os.path.join(self.go_path, "bin")

    def project_dir(self, *paths):
        return os.path.abspath(os.path.join(self.base_dir, *paths))

    def __rm(self, path):
        if self.cmd.exists(path):
            log.info("Removing '{}'...".format(path))
            self.cmd.rmf(path)
        else:
            log.debug("Nothing to delete for '{}'".format(path))

    def __rm_files(self, level):
        if level >= 0:
            self.__rm(self.project_dir("smart-contracts", "build"))  # truffle deploy
            self.__rm(self.project_dir("test", "integration", "vagrant", "data"))
            self.__rm(self.project_dir("test", "integration", "src", ".pytest_cache"))
            self.__rm(self.project_dir("test", "integration", "src", "py", ".pytest_cache"))
            self.__rm(self.cmd.get_user_home(".sifnoded"))  # Probably needed for "--keyring-backend test"

            self.__rm(self.project_dir("deploy", "networks"))  # from running integration tests

            # Peggy/devenv/hardhat cleanup
            # For full clean, also: cd smart-contracts && rm -rf node_modules && npm install
            # TODO Difference between yarn vs. npm install?
            self.__rm_hardhat_compiled_files()
            self.__rm_run_env_files()

            # Additional cleanup (not neccessary to make it work)
            # self.cmd.rm(self.project_dir("smart-contracts/combined.log"))
            # self.cmd.rmdir(self.project_dir("test/integration/.pytest_cache"))
            # self.cmd,rm(self.project_dir("smart-contracts/.env"))
            # self.cmd.rmdir(self.project_dir("deploy/networks"))
            # self.cmd.rmdir(self.project_dir("smart-contracts/.openzeppelin"))

            # docker image rm tendermintdev/sdk-proto-gen (used by Makefile on peggy2, used for "buf" command to build go bindings from ABI)

            # rmdir ~/.cache/yarn
        if level >= 1:
            for file in ["sifnoded", "ebrelayer", "sifgen"]:
                self.__rm(os.path.join(self.go_bin_dir, file))
            self.__rm(self.project_dir("smart-contracts", "node_modules"))
            self.__rm(self.project_dir("test", "localnet", "node_modules"))

        if level >= 2:
            if self.cmd.exists(self.go_path):
                self.cmd.execst(["chmod", "-R", "+w", self.go_path])
                self.__rm(self.go_path)
            self.__rm(self.cmd.get_user_home("go"))
            self.__rm(self.cmd.get_user_home(".npm"))
            self.__rm(self.cmd.get_user_home(".npm-global"))
            self.__rm(self.cmd.get_user_home(".cache/yarn"))
            self.__rm_run_env_files()
            self.__rm(self.cmd.get_user_home(".sifnode-integration"))

            # Peggy2
            self.__rm_peggy2_compiled_go_stubs()
            self.__rm(self.project_dir("smart-contracts", ".hardhat-compile"))

            # Remove go dependencies and re-download them (GOPATH=~/go)
            # rm -rv ~/go
            # mkdir ~/go
            # cd $PROJECT_DIR && go get -v -t -d ./...

            # On future/peggy2 these files are also created:
            # .proto-gen
            # .run/
            # cmd/ebrelayer/contract/generated/artifacts/
            # docs/peggy/node_modules/
            # smart-contracts/.hardhat-compile
            # smart-contracts/env.json
            # smart-contracts/environment.json

    def yarn(self, args, cwd=None, env=None):
        return self.cmd.execst(["yarn"] + args, cwd=cwd, env=env, pipe=False)

    def npx(self, args: Sequence[str], env: Optional[Mapping[str, str]] = None, cwd: Optional[str] = None,
        pipe: bool = True
    ) -> ExecResult:
        # Typically we want any npx commands to inherit stdout and strerr
        return self.cmd.execst(["npx"] + list(args), env=env, cwd=cwd, pipe=pipe)

    def run_peggy2_js_tests(self):
        # See smart-contracts/TEST.md:
        # 1. start environment
        # 2. npx hardhat test test/devenv/test_lockburn.ts --network localhost
        pass

    # Top-level "make install" should build everything, such as after git clone. If it does not, it's a bug.
    # "Official" way is "make clean install"
    def make_all(self, output_dir: Optional[str] = None):
        env = None if output_dir is None else {"GOBIN": output_dir}
        self.cmd.execst(["make", "install"], cwd=project_dir(), pipe=False, env=env)

    # IntegrationEnvironment
    # TODO Merge
    def make_go_binaries(self):
        # make go binaries (TODO Makefile needs to be trimmed down, especially "find")
        # cd test/integration; BASEDIR=... make
        # (checks all *.go files and, runs make in $BASEDIR, touches sifnoded, removes ~/.sifnoded/localnet
        self.cmd.execst(["make"], cwd=project_dir("test", "integration"), env={"BASEDIR": project_dir()}, pipe=False)

    # From PeggyEnvironment
    # TODO Merge
    # Main Makefile requires GOBIN to be set to an absolute path. Compiled executables ebrelayer, sifgen and
    # sifnoded will be written there. The directory will be created if it doesn't exist yet.
    def make_go_binaries_2(self, feature_toggles: Optional[Iterable[str]] = None):
        # Original: cd smart-contracts; make -C .. install
        extra_env = {feature: "1" for feature in feature_toggles}
        self.cmd.execst(["make", "install"], cwd=project_dir(), pipe=False, env=extra_env)

    def install_smart_contracts_dependencies(self):
        self.cmd.execst(["make", "clean-smartcontracts"], cwd=self.smart_contracts_dir)  # = rm -rf build .openzeppelin
        # According to peggy2, the plan is to move from npm install to yarn, but there are some issues with yarn atm.
        # self.yarn(["install"], cwd=self.smart_contracts_dir)
        # self.cmd.execst(["npm", "install"], cwd=self.smart_contracts_dir, pipe=False)
        self.npm_install(self.smart_contracts_dir)

    def write_vagrantenv_sh(self, state_vars, data_dir, ethereum_websocket_address, chainnet):
        env = dict_merge(state_vars, {
            # For running test/integration/execute_integration_tests_against_*.sh
            "TEST_INTEGRATION_DIR": project_dir("test/integration"),
            "TEST_INTEGRATION_PY_DIR": project_dir("test/integration/src/py"),
            "SMART_CONTRACTS_DIR": self.smart_contracts_dir,
            "datadir": data_dir,  # Needed by test_rollback_chain.py that calls ganache_start.sh
            "GANACHE_KEYS_JSON": os.path.join(data_dir, "ganachekeys.json"),  # Needed by test_rollback_chain.py that calls ganache_start.sh
            "ETHEREUM_WEBSOCKET_ADDRESS": ethereum_websocket_address,   # Needed by test_ebrelayer_replay.py (and possibly others)
            "CHAINNET": chainnet,   # Needed by test_ebrelayer_replay.py (and possibly others)
        })
        vagrantenv_path = project_dir("test/integration/vagrantenv.sh")
        self.cmd.write_text_file(vagrantenv_path, joinlines(format_as_shell_env_vars(env)))
        self.cmd.write_text_file(project_dir("test/integration/vagrantenv.json"), json.dumps(env))

    def get_peruser_config_dir(self):
        return self.cmd.get_user_home(".config", "siftool")

    def get_user_env_vars(self):
        env_file = os.environ["SIFTOOL_ENV_FILE"]
        return json.loads(self.cmd.read_text_file(env_file))

    def read_peruser_config_file(self, name):
        path = os.path.join(self.get_peruser_config_dir(), name + ".json")
        if self.cmd.exists(path):
            return json.loads(self.cmd.read_text_file(path))
        else:
            return None

    def init(self):
        self.clean()
        # self.cmd.rmdir(project_dir("smart-contracts/node_modules"))
        self.make_go_binaries_2()
        self.install_smart_contracts_dependencies()

    def __rm_sifnode_binaries(self):
        for filename in ["sifnoded", "ebrelayer", "sifgen"]:
            self.__rm(os.path.join(self.go_bin_dir, filename))

    # Removes hardhat-compiled smart contract files that are result of running
    # cd smart-contracts; npx hardhat run scripts/deploy_contracts.ts --network localhost
    def __rm_hardhat_compiled_files(self):
        for path in ["build", "artifacts", "cache", ".openzeppelin", ".hardhat-compile"]:
            self.__rm(os.path.join(self.smart_contracts_dir, path))

    def __rm_peggy2_compiled_go_stubs(self):
        # Peggy2: generated Go stubs (by smart-contracts/Makefile)
        if on_peggy2_branch:
            self.__rm(project_dir("cmd", "ebrelayer", "contract", "generated"))
        self.__rm(project_dir(".proto-gen"))

    def __rm_run_env_files(self):
        self.__rm(self.cmd.get_user_home(".sifnoded"))

        # Created by npx hardhat run scripts/devenv.ts and/or siftool run-env
        self.__rm(self.project_dir("smart-contracts", "relayerdb"))  # peggy1 only
        self.__rm(self.project_dir("test", "integration", "sifchainrelayerdb"))  # Probably obsolete on peggy2 TODO move to /tmp
        self.__rm(self.project_dir("smart-contracts", "environment.json"))
        self.__rm(self.project_dir("smart-contracts", "env.json"))
        self.__rm(self.project_dir("smart-contracts", ".env"))
        self.__rm(self.project_dir("smart-contracts", "venv"))
        self.__rm(self.project_dir(".run"))

    def clean(self):
        self.cmd.rmf(self.project_dir("smart-contracts", "node_modules"))
        self.cmd.rmf(os.path.join(self.siftool_dir, "build"))
        if on_peggy2_branch:
            for file in [".proto-gen", ".run", "cmd/ebrelayer/contract/generated/artifacts", "smart-contracts/.hardhat-compile"]:
                self.cmd.rmf(self.project_dir(file))
        else:
            # Output from "truffle compile" / "npx hardhat compile".
            # Wrong contents can cause hardhat to fail compilation after switching branches.
            self.cmd.rmf(self.project_dir("smart-contracts", "build"))
            self.cmd.rmf(self.project_dir("smart-contracts", "cache"))
            self.cmd.rmf(self.project_dir("smart-contracts", "artifacts"))

            self.__rm_sifnode_binaries()

    # Use this between run-env.
    def old_clean(self, level=None):
        level = 0 if level is None else int(level)
        force_kill_processes(self.cmd)
        self.__rm_files(level)

    def build(self):
        if on_peggy2_branch:
            self.npm_install(self.project_dir("smart-contracts"))
            self.npx(["hardhat", "compile"], cwd=self.project_dir("smart-contracts"), pipe=False)
        else:
            self.npm_install(self.project_dir("smart-contracts"))
            self.cmd.execst(["make", "install"], cwd=self.project_dir(), pipe=False)
            self.cmd.execst([self.project_dir("smart-contracts", "node_modules", ".bin", "truffle"), "compile"],
                cwd=self.project_dir("smart-contracts"), pipe=False)

    def rebuild(self):
        self.clean()
        self.build()

    def old_rebuild(self):
        # Use this after switching branches (i.e. develop vs. future/peggy2)
        self.clean(1)
        # self.cmd.execst(["npm", "install", "-g", "ganache-cli", "dotenv", "yarn"], cwd=self.smart_contracts_dir)
        self.install_smart_contracts_dependencies()
        self.cmd.execst(["make", "install"], cwd=self.project_dir(), pipe=False)

    def npm_install(self, path: str, disable_cache: bool = False):
        # TODO Add package-lock.json also on future/peggy2 branch?
        package_lock_json = os.path.join(path, "package.json" if on_peggy2_branch else "package-lock.json")
        sha1 = self.cmd.sha1_of_file(package_lock_json)
        node_modules = os.path.join(path, "node_modules")

        if self.cmd.exists(node_modules):
            cache_tag_file = os.path.join(node_modules, ".siftool-cache-tag")
            cache_tag = self.cmd.read_text_file(cache_tag_file) if self.cmd.exists(cache_tag_file) else None
            if (cache_tag is None) or (cache_tag != sha1):
                self.cmd.rmdir(node_modules)
            else:
                return

        assert not self.cmd.exists(node_modules)
        cache_dir = os.path.join(self.get_peruser_config_dir(), "npm-cache")
        cache_index = os.path.join(cache_dir, "index.json")
        cache = []
        if not self.cmd.exists(cache_dir):
            self.cmd.mkdir(cache_dir)
        if self.cmd.exists(cache_index):
            cache = json.loads(self.cmd.read_text_file(cache_index))
        idx = None
        for i, s in enumerate(cache):
            if s == sha1:
                idx = i
                break
        tarfile = os.path.join(cache_dir, "{}.tar".format(sha1))
        if idx is None:
            saved = dict(((f, self.cmd.read_text_file(f))
                for f in [os.path.join(path, x) for x in ["package-lock.json", "yarn.lock"]] if self.cmd.exists(f)))
            self.cmd.execst(["npm", "install"], cwd=path, pipe=False)
            cache_tag_file = os.path.join(node_modules, ".siftool-cache-tag")
            self.cmd.write_text_file(cache_tag_file, sha1)
            for file, contents in saved.items():
                self.cmd.write_text_file(file, contents)
            self.cmd.tar_create(node_modules, tarfile)
        else:
            cache.pop(idx)
            self.cmd.tar_extract(tarfile, node_modules)
        cache.insert(0, sha1)
        max_cache_items = 5
        if len(cache) > max_cache_items:
            for s in cache[max_cache_items:]:
                self.cmd.rm(os.path.join(cache_dir, "{}.tar".format(s)))
            cache = cache[:max_cache_items]
        self.cmd.write_text_file(cache_index, json.dumps(cache))

    def get_project_venv_dir(self):
        return project_dir("test", "integration", "framework", "venv")

    def project_python(self):
        return os.path.join(self.get_project_venv_dir(), "bin", "python3")

    def _ensure_build_dirs(self):
        for d in ["build", "build/repos", "build/generated"]:
            self.cmd.mkdir(os.path.join(self.siftool_dir, d))

    def generate_python_protobuf_stubs(self):
        # https://grpc.io/
        # https://grpc.github.io/grpc/python/grpc_asyncio.html
        self._ensure_build_dirs()
        project_proto_dir = self.project_dir("proto")
        third_party_proto_dir = self.project_dir("third_party", "proto")
        generated_dir = os.path.join(self.siftool_dir, "build/generated")
        repos_dir = os.path.join(self.siftool_dir, "build/repos")
        self.cmd.rmf(generated_dir)
        self.cmd.mkdir(generated_dir)
        cosmos_sdk_repo_dir = os.path.join(repos_dir, "cosmos-sdk")
        cosmos_proto_repo_dir = os.path.join(repos_dir, "cosmos-proto")
        # self.git_clone("https://github.com/gogo/protobuf", gogo_proto_dir, shallow=True)
        self.git_clone("https://github.com/cosmos/cosmos-sdk.git", cosmos_sdk_repo_dir, checkout_commit="dd65ef87322baa2023f195635890a2128a03d318")
        self.git_clone("https://github.com/cosmos/cosmos-proto.git", cosmos_proto_repo_dir, checkout_commit="213b76899fac883ac122728f7ab258166137be29")
        cosmos_sdk_proto_dir = os.path.join(cosmos_sdk_repo_dir, "proto")
        cosmos_proto_proto_dir = os.path.join(cosmos_proto_repo_dir, "proto")
        includes = [
            project_proto_dir,
            third_party_proto_dir,
            cosmos_sdk_proto_dir,
            cosmos_proto_proto_dir,
        ]

        # We cannot compile all proto files due to conflicting/inconsistent definitions (e.g. coin.proto).
        #
        # def find_proto_files(path, excludes=()):
        #     import re
        #     tmp = [os.path.relpath(i, start=path) for i in
        #         self.cmd.find_files(path, filter=lambda x: re.match(os.path.basename(x), "^(.*)\.proto$"))
        #     return sorted(list(set(tmp).difference(set(excludes) if excludes else set())))
        #
        # project_proto_files = find_proto_files(project_proto_dir)
        # third_party_proto_files = find_proto_files(third_party_proto_dir, excludes=[
        #     "cosmos/base/coin.proto",
        # ])
        # cosmos_sdk_proto_files = find_proto_files(cosmos_sdk_proto_dir, excludes=[
        #     "cosmos/base/query/v1beta1/pagination.proto",
        # ])
        # cosmos_proto_proto_files = find_proto_files(cosmos_proto_proto_dir)
        # proto_files = project_proto_files + third_party_proto_files + cosmos_sdk_proto_files + cosmos_proto_proto_files

        proto_files = [
            os.path.join(project_proto_dir, "sifnode/ethbridge/v1/tx.proto"),
            os.path.join(project_proto_dir, "sifnode/ethbridge/v1/query.proto"),
            os.path.join(project_proto_dir, "sifnode/ethbridge/v1/types.proto"),
            os.path.join(project_proto_dir, "sifnode/oracle/v1/network_descriptor.proto"),
            os.path.join(project_proto_dir, "sifnode/oracle/v1/types.proto"),
            os.path.join(third_party_proto_dir, "gogoproto/gogo.proto"),
            os.path.join(third_party_proto_dir, "google/api/annotations.proto"),
            os.path.join(third_party_proto_dir, "google/api/http.proto"),
            os.path.join(third_party_proto_dir, "cosmos/base/query/v1beta1/pagination.proto"),
            os.path.join(cosmos_sdk_proto_dir, "cosmos/tx/v1beta1/service.proto"),
            os.path.join(cosmos_sdk_proto_dir, "cosmos/base/abci/v1beta1/abci.proto"),
            os.path.join(cosmos_sdk_proto_dir, "cosmos/tx/v1beta1/tx.proto"),
            os.path.join(cosmos_sdk_proto_dir, "cosmos/tx/signing/v1beta1/signing.proto"),
            os.path.join(cosmos_sdk_proto_dir, "cosmos/crypto/multisig/v1beta1/multisig.proto"),
            os.path.join(cosmos_sdk_proto_dir, "cosmos/base/v1beta1/coin.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/abci/types.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/crypto/proof.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/crypto/keys.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/types/types.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/types/validator.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/types/params.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/types/block.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/types/evidence.proto"),
            os.path.join(cosmos_sdk_proto_dir, "tendermint/version/types.proto"),
            os.path.join(cosmos_proto_proto_dir, "cosmos_proto/cosmos.proto"),
        ]

        args = [self.project_python(), "-m", "grpc_tools.protoc"] + flatten(["-I", i] for i in includes) + [
            "--python_out", generated_dir, "--grpc_python_out", generated_dir] + proto_files
        self.cmd.execst(args, pipe=True)

    def git_clone(self, url, path, checkout_commit=None, shallow=False):
        if self.cmd.exists(os.path.join(path, ".git")):
            return
        log.debug("Cloning repository '{}' into '{}',,,".format(url, path))
        self.cmd.execst(["git", "clone", "-q"] + (["--depth", "1"] if shallow else []) + [url, path])
        if checkout_commit:
            self.cmd.execst(["git", "checkout", checkout_commit], cwd=path)

    def reset(self):
        self.__rm(os.path.join(self.smart_contracts_dir, "node_modules"))
        self.__rm_hardhat_compiled_files()
        self.__rm_sifnode_binaries()
        self.__rm(os.path.join(self.cmd.get_user_home(), ".sifnoded"))
        self.__rm_peggy2_compiled_go_stubs()
        self.__rm_run_env_files()
        self.npm_install(self.smart_contracts_dir)
        self.make_go_binaries_2()

    # Convenience wrapper
    def test(self, path: str, function: Optional[str] = None):
        test_arg = os.path.realpath(path)
        if function:
            test_arg += "::" + function
        args = [self.project_python(), "-m", "pytest", "-olog_cli=true", "-olog_cli_level=DEBUG", test_arg]
        self.cmd.execst(args, pipe=False)

    def pkill(self):
        for proc_name in ["node", "ebrelayer", "sifnoded", "geth"]:
            self.cmd.execst(["pkill", "--signal", "SIGTERM", proc_name], check_exit=False)

    def clean_run_env_state(self):
        self.__rm_run_env_files()
