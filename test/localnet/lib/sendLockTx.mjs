import { $ } from "zx";
import { getAddress } from "./getAddress.mjs";

export async function sendLockTx({
  binary,
  name,
  smartContractsDirectory = "/sifnode/smart-contracts",
  ethereumPrivateKey,
  deploymentName,
  amount,
}) {
  const srcAddr = await getAddress({ binary, name });

  process.env.ETHEREUM_PRIVATE_KEY = ethereumPrivateKey;

  const result = await $`
yarn \
    -s \
    --cwd ${smartContractsDirectory} \
    integrationtest:sendLockTx \
    --sifchain_address ${srcAddr}
    --symbol eth \
    --ethereum_private_key_env_var \
    --json_path ${smartContractsDirectory}/deployments/${deploymentName} \
    --gas estimate \
    --ethereum_network ropsten \
    --bridgebank_address ${bridgeBankAddress} \
    --ethereum_address ${ethereumAddress} \
    --amount ${amount}
    `;

  return result;
}
