import os
import shutil
from git import Repo
import sys
import subprocess
import json
from mnemonic import Mnemonic
import time
import toml
import socket
from faker import Faker
from siftool import command, project, common

#definite variables.
home_directory_path = os.environ["HOME"]
if not home_directory_path.endswith("/"):
    home_directory_path += "/"
go_path = home_directory_path + "go/bin"
sifnode_repo_url = "https://github.com/Sifchain/sifnode.git"
branch = "develop"
chain_id = "local-sifnode-1"
number_of_nodes = 5
sifnode_log_level = "info"

#pip install gitpython
#pip install mnemonic
#pip install requests

cmd = command.Command()
prj = project.Project(cmd, common.project_dir())

def background_cli_command(command, gopath):
    command_eddition = "export PATH={go_path}:$PATH && ".format(go_path=gopath)
    command = command_eddition + command
    print(command)
    subprocess.Popen(command, shell=True)

def cli_command(command, gopath):
    command_eddition = "export PATH={go_path}:$PATH &&".format(go_path=gopath)
    command = command_eddition + command
    try:
        output = subprocess.check_output(command, stderr=subprocess.STDOUT, shell=True, universal_newlines=True)
    except subprocess.CalledProcessError as exc:
        print("Status : FAIL", exc.returncode, exc.output)
        return False
    else:
        return output

def clean_up_old_installation(home_directory_path):
    home_directory_list = os.listdir(home_directory_path)
    for directory in home_directory_list:
        if ".sifnode" in directory:
            shutil.rmtree(home_directory_path + directory)
    # for dirname in [".sifnode"]:
    #     path = os.path.join(home_directory_path, dirname)
    #     if cmd.exists(path):
    #         shutil.rmtree(path)

def clean_up_old_clone(directory_path):
    try:
        shutil.rmtree(directory_path)
    except:
        print("Not Cloned")

def clone_sifnode_repo(sifnode_git_url):
    repo = Repo.clone_from(sifnode_git_url, "./sifnoded")
    return repo

def generate_key(key_name, sif_folder, mnemonic, home_directory_path):
    print("Check if key exists.")
    result = cli_command(
        "sifnoded keys list --keyring-backend test --output json --home {home_directory_path}.{sif_folder}".format(
            home_directory_path=home_directory_path, sif_folder=sif_folder), go_path)
    if result != '[]':
        result_json_object = json.loads(result)

    print("Generate Sif Account Key")
    command = 'echo "{mnemonic}" | sifnoded keys add {key_name} --recover --keyring-backend test --home {home_directory_path}.{sif_folder}'.format(
        home_directory_path=home_directory_path, sif_folder=sif_folder, key_name=key_name, mnemonic=mnemonic)
    print(command)
    key_result = cli_command(command, go_path)
    if key_name in key_result:
        print("Key Created")
    else:
        print("Key Not Created")

    result = cli_command(
        "sifnoded keys list --keyring-backend test --output json --home {home_directory_path}.{sif_folder}".format(
            home_directory_path=home_directory_path, sif_folder=sif_folder), go_path)

    result_json_object = json.loads(result)
    key_load_success = False
    sif_account_address = ""

    for key in result_json_object:
        if key["name"] == key_name:
            key_load_success = True
            sif_account_address = key["address"]
    if key_load_success and sif_account_address:
        print("Key Loaded")
        return sif_account_address
    else:
        print("Key Not Loaded")
        sys.exit(1)

def init_chain(sif_folder, chain_id):
    command = "sifnoded init test --chain-id={chain_id} -o --home {home_directory_path}.{sif_folder}".format(
        home_directory_path=home_directory_path, chain_id=chain_id, sif_folder=sif_folder)
    print(command)
    cli_command(command, go_path)
    directory_to_check = "{home_directory_path}.{sif_folder}".format(home_directory_path=home_directory_path, sif_folder=sif_folder)
    print("Validate Init, check to see if init directory exsits.", directory_to_check)
    if os.path.isdir(directory_to_check):
        print("directory_exists")
    else:
        sys.exit(1)

def add_multisig_key(home_directory_path, home_folder_directory):
    command = "sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test --home {home_directory_path}.{home_folder_directory}".format(
        home_directory_path=home_directory_path, home_folder_directory=home_folder_directory)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if "mkey" in result:
        print("Multisig Key Added Succesfully for Akasha/Sif")
    else:
        print("Multisig Key was unable to be added for Akasha/Sif")
        sys.exit(1)

def add_genesis_account(key_name, home_folder_directory, home_folder):
    command = "sifnoded add-genesis-account {key_name} 990000000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home {home_folder_directory}.{home_folder}".format(key_name=key_name, home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def add_genesis_clp_account(key_name, home_folder_directory, home_folder):
    command = "sifnoded add-genesis-clp-admin {key_name} --keyring-backend=test --home {home_folder_directory}.{home_folder}".format(key_name=key_name, home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def set_whitelist_admin(key_name, home_folder_directory, home_folder):
    command = "sifnoded set-genesis-whitelister-admin {key_name} --keyring-backend=test --home {home_folder_directory}.{home_folder}".format(key_name=key_name, home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def set_whitelist(home_folder_directory, home_folder, whitelist_denoms_json_path):
    command = "sifnoded set-gen-denom-whitelist {whitelist_denoms_json_path} --home {home_folder_directory}.{home_folder}".format(whitelist_denoms_json_path=whitelist_denoms_json_path, home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def set_genesis_validator(key_name, home_folder_directory, home_folder):
    command = "sifnoded add-genesis-validators $(sifnoded keys show {key_name} -a --bech val --keyring-backend=test --home {home_folder_directory}.{home_folder}) --keyring-backend=test --home {home_folder_directory}.{home_folder}".format(key_name=key_name, home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def gen_tx(key_name, home_folder_directory, home_folder, chain_id):
    command = "sifnoded gentx {key_name} --keyring-backend=test 1000000000000000000000000rowan --chain-id={chain_id} --home {home_folder_directory}.{home_folder}".format(chain_id=chain_id, key_name=key_name, home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def collect_gentx(home_folder_directory, home_folder):
    command = "sifnoded collect-gentxs --home {home_folder_directory}.{home_folder}".format(home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

def validate_genesis(home_folder_directory, home_folder):
    command = "sifnoded validate-genesis --home {home_folder_directory}.{home_folder}".format(home_folder_directory=home_folder_directory, home_folder=home_folder)
    print(command)
    result = cli_command(command, go_path)
    print(result)
    if result == False:
        print("failure")
        return False
    else:
        return True

faker = Faker()

print("Init Mnemonic")
mnemo = Mnemonic("english")

account_mnemonics = {}
for node_number in range(number_of_nodes):
    if node_number == 0:
        account_mnemonics["0"] = {"name":"sif","mnemonic": "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"}
        account_mnemonics["0_a"] = {"name":"akasha","mnemonic": "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard"}
    else:
        words = mnemo.generate(strength=256)
        account_mnemonics[str(node_number)] = {"name":"sif-"+str(node_number),"mnemonic": words}

print("Cleanup old installs of sifnoded.")
clean_up_old_installation(home_directory_path)

# print("Cleanup old clone of sifnoded repo")
# clean_up_old_clone("./sifnoded")
#
# print("Clone fresh sifnoded repo.")
# repo = clone_sifnode_repo(sifnode_repo_url)
#
# print("Checkout branch for sifnoded repo.")
# repo.git.checkout(branch)
#
# print("Build sifnoded for the branch specified.")
# result = cli_command("cd sifnoded/ && make clean install", go_path)
# print(result)
#
# print("Check sifnoded installed")
# if not cli_command("sifnoded keys list --keyring-backend test --output json", go_path):
#     print("Sifnode not installed")
# else:
#     print("Successfully Installed")

print("Loop and setup local configs for the number of desired nodes.")
first = True
for node_number in range(number_of_nodes):

    sifnode_home_directory_name = "sifnode-" + str(node_number)
    p2p_port = int("276" + str(node_number))
    grpc_port = int("909" + str(node_number))
    grpc_web_port = int("919" + str(node_number))
    address_port = int("276" + str(node_number))
    rpc_port = int("286" + str(node_number))
    api_port = int("131" + str(node_number))
    pprof_port = int("606" + str(node_number))

    account_mnemonics[str(node_number)]["home_directory"] = "sifnode-" + str(node_number)
    account_mnemonics[str(node_number)]["p2p_port"] = p2p_port
    account_mnemonics[str(node_number)]["grpc_port"] = grpc_port
    account_mnemonics[str(node_number)]["grpc_web_port"] = grpc_web_port
    account_mnemonics[str(node_number)]["address_port"] = address_port
    account_mnemonics[str(node_number)]["rpc_port"] = rpc_port
    account_mnemonics[str(node_number)]["api_port"] = api_port
    account_mnemonics[str(node_number)]["pprof_port"] = pprof_port

    print("Home Directory: ~/." + sifnode_home_directory_name)
    print("P2P Port: ", p2p_port)
    print("GRPC Port: ", grpc_port)
    print("GRPC Web Port: ", grpc_web_port)
    print("ADDRESS Port: ", address_port)
    print("RPC Port: ", rpc_port)
    print("API Port: ", api_port)
    print("ChainID: ", chain_id)

    if first:
        # For the first validator we setup the genesis
        print("Init Chain")

        init_chain(sifnode_home_directory_name, chain_id)

        for account in account_mnemonics:
            if account == "0":
                print("Generate Sif Key")
                sif_account_address = generate_key(account_mnemonics["0"]["name"], sifnode_home_directory_name, account_mnemonics["0"]["mnemonic"],home_directory_path)
                account_mnemonics["0"]["account_address"] = sif_account_address
            elif account == "0_a":
                print("Generate Akasha Key")
                akasha_account_address = generate_key(account_mnemonics["0_a"]["name"], sifnode_home_directory_name, account_mnemonics["0_a"]["mnemonic"],home_directory_path)
                account_mnemonics["0_a"]["account_address"] = akasha_account_address
            else:
                dynamic_gen_account_address = generate_key(account_mnemonics[account]["name"], sifnode_home_directory_name, account_mnemonics[account]["mnemonic"],home_directory_path)
                account_mnemonics[account]["account_address"] = dynamic_gen_account_address

        print("Generate Multi Sif Key")
        add_multisig_key(home_directory_path,sifnode_home_directory_name)

        print("Generate Genesis Account for Account:", sif_account_address)
        if not add_genesis_account(sif_account_address, home_directory_path, sifnode_home_directory_name):
            print("Failed to Add Genesis Account")
            sys.exit(1)

        print("Generate Genesis Account for Account:", akasha_account_address)
        if not add_genesis_account(akasha_account_address, home_directory_path, sifnode_home_directory_name):
            print("Failed to Add Genesis Account")
            sys.exit(1)

        for account in account_mnemonics:
            if account == "0" or account == "0_a":
                continue
            else:
                print("Generate Genesis Account for Account:", account_mnemonics[account]["account_address"])
                if not add_genesis_account(account_mnemonics[account]["account_address"], home_directory_path, sifnode_home_directory_name):
                    print("Failed to Add Genesis Account")
                    sys.exit(1)

        print("Generate Genesis Account CLP Admin for Account:", sif_account_address)
        if not add_genesis_clp_account(sif_account_address, home_directory_path, sifnode_home_directory_name):
            print("Failed to Add Genesis CLP Admin Account")
            sys.exit(1)

        print("Generate Genesis Account CLP Admin Account:", akasha_account_address)
        if not add_genesis_clp_account(akasha_account_address, home_directory_path, sifnode_home_directory_name):
            print("Failed to Add Genesis Account")
            sys.exit(1)

        print("Generate Genesis Account for Account:", sif_account_address)
        if not set_whitelist_admin(sif_account_address, home_directory_path, sifnode_home_directory_name):
            print("Failed to Add Whitelist Admin Account")
            sys.exit(1)

        print("Set Whitelist from sifnode/scripts/denoms.json")
        denoms_json_path = prj.project_dir("scripts/denoms.json")
        if not set_whitelist(home_directory_path, sifnode_home_directory_name, denoms_json_path):
            print("Whitelist failed to be set.")
            sys.exit(1)

        print("Generate Genesis Validator for Account:", sif_account_address)
        if not set_genesis_validator("sif", home_directory_path, sifnode_home_directory_name):
            print("Failed to Generate Genesis Validator Account")
            sys.exit(1)

        for account in account_mnemonics:
            if account == "0" or account == "0_a":
                continue
            else:
                print("Set Genesis Validator:", account_mnemonics[account]["account_address"])
                if not set_genesis_validator("sif-" + str(account), home_directory_path, sifnode_home_directory_name):
                    print("Failed to Generate Genesis Validator Account")
                    sys.exit(1)

        print("GenTX:", sif_account_address)
        if not gen_tx("sif", home_directory_path, sifnode_home_directory_name, chain_id):
            print("Failed to GENTX")
            sys.exit(1)

        print("Collect GenTX:", sif_account_address)
        if not collect_gentx(home_directory_path, sifnode_home_directory_name):
            print("Failed to Collect GENTX")
            sys.exit(1)

        print("Validate Genesis:", sif_account_address)
        if not validate_genesis(home_directory_path, sifnode_home_directory_name):
            print("Failed to Validate Genesis")
            sys.exit(1)

        print("Get NodeID")
        node_id = cli_command("sifnoded tendermint show-node-id --home "+home_directory_path + "." + sifnode_home_directory_name, go_path)
        account_mnemonics[str(node_number)]["node_id"] = node_id.strip()

        #EDIT GENESIS HERE
        open_genesis_file = open(home_directory_path + "." + sifnode_home_directory_name + "/config/genesis.json", "r").read()
        genesis_json_object = json.loads(open_genesis_file)
        genesis_json_object["app_state"]["gov"]["voting_params"] = {"voting_period": "120s"}
        genesis_json_object["app_state"]["crisis"]["constant_fee"] = {"denom": "rowan", "amount": "1000"}
        genesis_json_object["app_state"]["gov"]["deposit_params"]["min_deposit"] = [{"denom": "rowan", "amount": "10000000"}]
        genesis_json_object["app_state"]["staking"]["params"]["bond_denom"] = "rowan"
        genesis_json_object["app_state"]["mint"]["params"]["mint_denom"] = "rowan"


        genesis_payload = json.dumps(genesis_json_object)

        print("Edited Genesis Payload")
        print(genesis_payload)
        open_genesis_file = open(home_directory_path + "." + sifnode_home_directory_name + "/config/genesis.json","w")
        open_genesis_file.write(genesis_payload)
        open_genesis_file.close()
        print("---------NODE {node_number} CONFIGURED----------".format(node_number=str(node_number)))
        print("----\n\n\n\n\n")
        first = False
    else:
        # For every other validateor we call "sifnoded init" and overwrite the genesis file with the one from the first validator
        print("Init Node: ", node_number)
        init_chain(sifnode_home_directory_name, chain_id)

        print("Setup Genesis File Location Variables")
        genesis_file_location = home_directory_path + ".sifnode-0/config/genesis.json"
        genesis_file_destination = home_directory_path + "."+sifnode_home_directory_name + "/config/genesis.json"

        print("Remove Init Genesis File")
        os.remove(genesis_file_destination)

        print("Copy Genesis File from Node-0")
        shutil.copyfile(genesis_file_location,genesis_file_destination)

        # For every other validateor we call "sifnoded init" and overwrite the genesis file with the one from the first validator
        print("Load Node Account")
        dynamic_gen_account_address = generate_key(account_mnemonics[str(node_number)]["name"], sifnode_home_directory_name, account_mnemonics[str(node_number)]["mnemonic"],home_directory_path)
        account_mnemonics[str(node_number)]["account_address"] = dynamic_gen_account_address

        print("Get NodeID")
        node_id = cli_command("sifnoded tendermint show-node-id --home "+home_directory_path + "." + sifnode_home_directory_name, go_path)
        account_mnemonics[str(node_number)]["node_id"] = node_id.strip()
        print("---------NODE {node_number} CONFIGURED----------".format(node_number=str(node_number)))
        print("----\n\n\n\n\n")

print("Loop generated account dictionary and create list of peers, and start nodes once the tomls have been edited. Also we will build a report output for the user here.")
peers = ""
for account in account_mnemonics:
    if account == "0":
        seed = account_mnemonics[account]["node_id"] + "@127.0.0.1:" + str(account_mnemonics[account]["p2p_port"])
        peers += account_mnemonics[account]["node_id"] + "@127.0.0.1:" + str(account_mnemonics[account]["p2p_port"])
        account_mnemonics[account]["seed"] = seed
    elif account == "0_a":
        continue

print("Peers: ", peers)

node_information = ""

for account in account_mnemonics:
    if account != "0_a":
        app_toml_location = home_directory_path + "." + "sifnode-" + str(account) + "/config/app.toml"
        config_toml_location = home_directory_path + "." + "sifnode-" + str(account) + "/config/config.toml"
        sifnode_home_directory_name = "sifnode-" + str(account)
        p2p_port = str(account_mnemonics[account]["p2p_port"])
        grpc_port = str(account_mnemonics[account]["grpc_port"])
        grpc_web_port = str(account_mnemonics[account]["grpc_web_port"])
        address_port = str(account_mnemonics[account]["address_port"])
        rpc_port = str(account_mnemonics[account]["rpc_port"])
        api_port = str(account_mnemonics[account]["api_port"])

    if account == "0":
        print("Edit sifnode/config/app.toml minimum-gas-prices and api enabled")
        data = toml.load(app_toml_location)
        data["minimum-gas-prices"] = "0.5rowan"
        data['api']['enable'] = True
        api_address = "tcp://0.0.0.0:" + str(account_mnemonics[account]["api_port"])
        data["api"]["address"] = api_address
        f = open(app_toml_location, 'w')
        toml.dump(data, f)
        f.close()

        print("Edit sifnode/config/config.toml peers, seeds, external_address")
        data = toml.load(config_toml_location)
        hostname = socket.gethostname()
        ip_address = socket.gethostbyname(hostname)

        data['log_level'] = sifnode_log_level
        data['p2p']["external_address"] = str(ip_address) + ":" + str(account_mnemonics[account]["p2p_port"])
        data['p2p']['max_num_inbound_peers'] = 50
        data['p2p']['max_num_outbound_peers'] = 50
        data['p2p']['allow_duplicate_ip'] = True
        data["rpc"]["pprof_laddr"] = "127.0.0.1:" + str(account_mnemonics[str(node_number)]["pprof_port"])

        data['moniker'] = str(faker.first_name())
        account_mnemonics[account]["moniker"] = data['moniker']

        f = open(config_toml_location, 'w')
        toml.dump(data, f)
        f.close()

        print("Start Sifnode-"+ str(account))
        background_cli_command("""sifnoded start --home {home_folder_directory}.{home_folder} \\
            --p2p.laddr 127.0.0.1:{p2p_port}  \\
            --grpc.address 0.0.0.0:{grpc_port} \\
            --grpc-web.address 0.0.0.0:{grpc_web_port} \\
            --address tcp://0.0.0.0:{address_port} \\
            --rpc.laddr tcp://127.0.0.1:{rpc_port}""".format(p2p_port=p2p_port,grpc_port=grpc_port,grpc_web_port=grpc_web_port,address_port=address_port,rpc_port=rpc_port,home_folder_directory=home_directory_path, home_folder=sifnode_home_directory_name), go_path)

        node_information += """
        ---------------------------
        Node # {node_number}
        Commad: export PATH={go_path}:$PATH && sifnoded start --home {home_folder_directory}.{home_folder} --p2p.laddr 127.0.0.1:{p2p_port} --grpc.address 0.0.0.0:{grpc_port} --grpc-web.address 0.0.0.0:{grpc_web_port} --address tcp://0.0.0.0:{address_port} --rpc.laddr tcp://127.0.0.1:{rpc_port}
        Moniker:  {moniker}
        p2p Port: {p2p_port} 
        grpc Port: {grpc_port}
        grpc web Port: {grpc_web_port}
        address Port: {address_port}
        rpc Port: {rpc_port}
        api Port: {api_port}
        peers: {peers}
        ---------------------------
        """.format(peers=peers,
                   node_number=str(account),
                   api_port=api_port,
                   moniker=data['moniker'],
                   p2p_port=p2p_port,
                   grpc_port=grpc_port,
                   grpc_web_port=grpc_web_port,
                   address_port=address_port,
                   rpc_port=rpc_port,
                   home_folder_directory=home_directory_path,
                   home_folder=sifnode_home_directory_name,
                   go_path=go_path)

    elif account == "0_a":
        continue
    else:
        print("sleep for 10 seconds..")
        time.sleep(10)
        print("Edit sifnode/config/app.toml minimum-gas-prices and api enabled")
        data = toml.load(app_toml_location)
        data["minimum-gas-prices"] = "0.5rowan"
        data['api']['enable'] = True
        api_address = "tcp://0.0.0.0:" + str(account_mnemonics[account]["api_port"])
        data["api"]["address"] = api_address
        f = open(app_toml_location, 'w')
        toml.dump(data, f)
        f.close()

        print("Edit sifnode/config/config.toml peers, seeds, external_address")
        data = toml.load(config_toml_location)
        hostname = socket.gethostname()
        ip_address = socket.gethostbyname(hostname)

        data['log_level'] = sifnode_log_level
        data['p2p']["external_address"] = str(ip_address) + ":" + str(account_mnemonics[account]["p2p_port"])
        data['p2p']['persistent_peers'] = peers
        data['p2p']['max_num_inbound_peers'] = 50
        data['p2p']['max_num_outbound_peers'] = 50
        data['p2p']['max_num_outbound_peers'] = 50
        data['p2p']['allow_duplicate_ip'] = True
        data["rpc"]["pprof_laddr"] = "127.0.0.1:" + str(account_mnemonics[str(node_number)]["pprof_port"])

        data['moniker'] = str(faker.first_name())
        account_mnemonics[account]["moniker"] = data['moniker']

        f = open(config_toml_location, 'w')
        toml.dump(data, f)
        f.close()

        staking_command = """
        sifnoded tx staking create-validator \\
        --amount=1000000000000000000000000rowan \\
        --pubkey=$(sifnoded tendermint show-validator --home {home_folder_directory}.{home_folder}) \\
        --moniker={moniker} \\
        --chain-id={chain_id} \\
        --commission-rate="0.10" \\
        --commission-max-rate="0.20" \\
        --commission-max-change-rate="0.01" \\
        --min-self-delegation="1000000" \\
        --node=tcp://127.0.0.1:2860 \\
        --fees 1000000000000000000rowan \\
        --keyring-backend test \\
        --from=sif-{node_number} \\
        -y \\
        --home {home_folder_directory}.{home_folder}""".format(home_folder_directory=home_directory_path, home_folder=sifnode_home_directory_name, chain_id=chain_id,moniker=account_mnemonics[account]["moniker"], node_number=str(account))
        print(staking_command)
        command_output = cli_command(staking_command, go_path)
        print(command_output)
        if not command_output:
            print("Staking comman failed.")

        print("Start Sifnode-"+ str(account))
        background_cli_command("""sifnoded start --home {home_folder_directory}.{home_folder} \\
                --p2p.laddr 127.0.0.1:{p2p_port}  \\
                --grpc.address 0.0.0.0:{grpc_port} \\
                --grpc-web.address 0.0.0.0:{grpc_web_port} \\
                --address tcp://0.0.0.0:{address_port} \\
                --rpc.laddr tcp://127.0.0.1:{rpc_port}""".format(p2p_port=p2p_port,grpc_port=grpc_port,grpc_web_port=grpc_web_port,address_port=address_port,rpc_port=rpc_port,home_folder_directory=home_directory_path, home_folder=sifnode_home_directory_name), go_path)

        node_information += """
        ---------------------------
        Node # {node_number}
        Commad: export PATH={go_path}:$PATH && sifnoded start --home {home_folder_directory}.{home_folder} --p2p.laddr 127.0.0.1:{p2p_port} --grpc.address 0.0.0.0:{grpc_port} --grpc-web.address 0.0.0.0:{grpc_web_port} --address tcp://0.0.0.0:{address_port} --rpc.laddr tcp://127.0.0.1:{rpc_port}
        Moniker:  {moniker}
        p2p Port: {p2p_port} 
        grpc Port: {grpc_port}
        grpc web Port: {grpc_web_port}
        address Port: {address_port}
        rpc Port: {rpc_port}
        api Port: {api_port}
        peers: {peers}
        ---------------------------
        """.format(peers=peers,
                   node_number=str(account),
                   api_port=api_port,
                   moniker=data['moniker'],
                   p2p_port=p2p_port,
                   grpc_port=grpc_port,
                   grpc_web_port=grpc_web_port,
                   address_port=address_port,
                   rpc_port=rpc_port,
                   home_folder_directory=home_directory_path,
                   home_folder=sifnode_home_directory_name,
                   go_path=go_path)


node_information_log = open("node_information_log.log", "w")
print(node_information)
node_information_log.write(node_information)
node_information_log.close()

while True:
    print("Keeping Process Alive to Keep the Sifnode Simulator Alive, if you want it to run without this script running simply open node_information_log.log and run the start commands for each node in the log.")
    time.sleep(10)
