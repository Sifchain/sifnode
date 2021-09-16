import { GolangResults } from "./golangBuilder";
import { SifnodedResults, SifnodedRunner, ValidatorValues } from "./sifnoded";
import { SmartContractDeployResult } from "./smartcontractDeployer";
import { EthereumResults } from "./devEnv";
import path from 'path';
import fs from 'fs';

interface ETHEnv {
  ETH_CHAIN_ID: number,
  ETH_HOST: string,
  ETH_PORT: number,
  ETHEREUM_ADDRESS: string,
  ETHEREUM_PRIVATE_KEY: string,
  ROWAN_SOURCE: string,
  ETH_ACCOUNT_OPERATOR_ADDRESS: string
  ETH_ACCOUNT_OPERATOR_PRIVATEKEY: string
  ETH_ACCOUNT_OWNER_ADDRESS: string
  ETH_ACCOUNT_OWNER_PRIVATEKEY: string
  ETH_ACCOUNT_PAUSER_ADDRESS: string
  ETH_ACCOUNT_PAUSER_PRIVATEKEY: string
  ETH_ACCOUNT_PROXYADMIN_ADDRESS: string
  ETH_ACCOUNT_PROXYADMIN_PRIVATEKEY: string
  ETH_ACCOUNT_VALIDATOR_ADDRESS: string
  ETH_ACCOUNT_VALIDATOR_PRIVATEKEY: string
}

interface ContractEnv {
  BRIDGE_BANK_ADDRESS: string
  BRIDGE_REGISTERY_ADDRESS: string
  COSMOS_BRIDGE_ADDRESS: string
  ROWANTOKEN_ADDRESS: string
  BRIDGE_TOKEN_ADDRESS: string // Same address as Rowantoken
}

interface GOEnv {
  GOBIN: string
}

interface SifEnv {
  TCP_URL: string
  VALIDATOR_MENOMONIC: string
  VALIDATOR_MONIKER: string
  VALIDATOR_PASSWORD: string
  VALIDATOR_PUB_KEY: string
  VALIDATOR_ADDRESS: string
  VALIDATOR_CONSENSUS_ADDRESS: string
  CHAINDIR: string
}
interface EnvOutput {
  COMPUTED: {
    BASEDIR: string
  }
  SIFNODE?: SifEnv
  GOLANG?: GOEnv
  CONTRACTS?: ContractEnv
  ETHEREUM?: ETHEnv
}

interface EnvDictionary {
  [key: string]: string
}

export function EnvJSONWriter(args: {
  ethResults?: EthereumResults,
  goResults?: GolangResults,
  sifResults?: SifnodedResults,
  contractResults?: SmartContractDeployResult,
}) {
  const baseDir = path.resolve(__dirname, "../../..")
  const output: EnvOutput = {
    COMPUTED: {
      BASEDIR: baseDir
    }
  };
  if (args.ethResults != undefined) {
    const eth = args.ethResults
    const env: ETHEnv = {
      ETHEREUM_ADDRESS: eth.accounts.available[0].address,
      ETHEREUM_PRIVATE_KEY: eth.accounts.available[0].privateKey,
      ROWAN_SOURCE: eth.accounts.operator.privateKey,
      ETH_ACCOUNT_OPERATOR_ADDRESS: eth.accounts.operator.address,
      ETH_ACCOUNT_OPERATOR_PRIVATEKEY: eth.accounts.operator.privateKey,
      ETH_ACCOUNT_OWNER_ADDRESS: eth.accounts.owner.address,
      ETH_ACCOUNT_OWNER_PRIVATEKEY: eth.accounts.owner.privateKey,
      ETH_ACCOUNT_PAUSER_ADDRESS: eth.accounts.pauser.address,
      ETH_ACCOUNT_PAUSER_PRIVATEKEY: eth.accounts.pauser.privateKey,
      ETH_ACCOUNT_PROXYADMIN_ADDRESS: eth.accounts.proxyAdmin.address,
      ETH_ACCOUNT_PROXYADMIN_PRIVATEKEY: eth.accounts.proxyAdmin.privateKey,
      ETH_ACCOUNT_VALIDATOR_ADDRESS: eth.accounts.validators[0].address,
      ETH_ACCOUNT_VALIDATOR_PRIVATEKEY: eth.accounts.validators[0].privateKey,
      ETH_CHAIN_ID: eth.chainId,
      ETH_HOST: eth.httpHost,
      ETH_PORT: eth.httpPort,
    }
    output.ETHEREUM = env
  }
  if (args.contractResults != undefined) {
    const contract = args.contractResults.contractAddresses
    const env: ContractEnv = {
      BRIDGE_BANK_ADDRESS: contract.bridgeBank,
      BRIDGE_REGISTERY_ADDRESS: contract.bridgeRegistry,
      COSMOS_BRIDGE_ADDRESS: contract.cosmosBridge,
      ROWANTOKEN_ADDRESS: contract.rowanContract,
      BRIDGE_TOKEN_ADDRESS: contract.rowanContract
    }
    output.CONTRACTS = env
  }
  if (args.goResults != undefined) {
    output.GOLANG = {
      GOBIN: args.goResults.goBin
    }
  }
  if (args.sifResults != undefined) {
    const sif = args.sifResults
    const val = sif.validatorValues[0]
    const env: SifEnv = {
      TCP_URL: sif.tcpurl,
      VALIDATOR_ADDRESS: val.address,
      VALIDATOR_CONSENSUS_ADDRESS: val.validator_consensus_address,
      VALIDATOR_MENOMONIC: val.mnemonic,
      VALIDATOR_MONIKER: val.moniker,
      VALIDATOR_PASSWORD: val.password,
      VALIDATOR_PUB_KEY: val.pub_key,
      // TODO: Remove hardcoded strings
      CHAINDIR: path.resolve("/tmp/sifnodedNetwork/validators", val.chain_id, val.moniker)
    }
    output.SIFNODE = env
  }
  try {
    const envValues: string[] = []
    const rootValues: EnvDictionary = {}
    Object.values(output).forEach(module => {
      Object.entries(module).forEach(entry => {
        envValues.push(`${entry[0]}="${entry[1]}"`);
        rootValues[entry[0]] = entry[1] as string;
      });
    });
    const envText = envValues.join("\n")
    fs.writeFileSync(path.resolve(__dirname, "../../", ".env"), envText);
    fs.writeFileSync(path.resolve(__dirname, "../../", "env.json"), JSON.stringify(rootValues))
    fs.writeFileSync(path.resolve(__dirname, "../../", "environment.json"), JSON.stringify(args));
    console.log("Wrote environment and JSON values to disk. PATH: ", path.resolve(__dirname));
  }
  catch (error) {
    console.error("Failed to write environment/json values to disk, ERROR: ", error);
  }
}