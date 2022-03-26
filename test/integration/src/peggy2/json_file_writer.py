import json
from siftool import run_env

cmd = run_env.Integrator()
file_path = cmd.project.project_dir("smart-contracts/src/devenv/ibc_registry.json")

def write_registry_json(denom, address):
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
      "network": "NETWORK_DESCRIPTOR_UNSPECIFIED",
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

write_registry_json("ibc", "address")