import json
from siftool import run_env

cmd = run_env.Integrator()
ibc_json_file_path = cmd.project.project_dir("/tmp/ibc_registry.json")
double_peggy_json_file_path = cmd.project.project_dir("/tmp/double_peggy_registry.json")

def write_registry_json(denom, address, isIbc=True):
    if isIbc is True:
        network = "NETWORK_DESCRIPTOR_UNSPECIFIED"
        file_path = ibc_json_file_path
    else:
        network = "NETWORK_DESCRIPTOR_ETHEREUM_TESTNET_ROPSTEN"
        file_path = double_peggy_json_file_path
    data = {
      "is_whitelisted": True,
      "decimals": "18",
      "denom": denom,
      "base_denom": denom,
      "path": "",
      "ibc_channel_id": "",
      "ibc_counterparty_channel_id": "",
      "display_name": "",
      "display_symbol": "",
      "network": network,
      "address": address,
      "external_symbol": "",
      "transfer_limit": "",
      "permissions": [],
      "unit_denom": "",
      "ibc_counterparty_denom": "",
      "ibc_counterparty_chain_id": ""
    }

    entries = {
      "entries": [data]
    }

    with open(file_path, "w") as outfile:
      outfile.write(json.dumps(entries, indent=4))

# write_registry_json("ibc", "address")