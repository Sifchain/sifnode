import createEthbridgeService from ".";
import { Asset, AssetAmount, Token } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { advanceBlock } from "../../test/utils/advanceBlock";
import { juniper, akasha, ethAccounts } from "../../test/utils/accounts";
import {
  createTestEthService,
  createTestSifService,
} from "../../test/utils/services";
import { getBalance, getTestingTokens } from "../../test/utils/getTestingToken";
import config from "../../config.localnet.json";
import Web3 from "web3";
import JSBI from "jsbi";
import { sleep } from "../../test/utils/sleep";

const [ETH, CETH, ATK, CATK, ROWAN, EROWAN] = getTestingTokens([
  "ETH",
  "CETH",
  "ATK",
  "CATK",
  "ROWAN",
  "EROWAN",
]);

async function waitFor(
  getter: () => Promise<any>,
  expected: any,
  name: string
) {
  console.log(`Starting wait: "${name}"`);
  let value: any;
  for (let i = 0; i < 100; i++) {
    await sleep(1000);
    value = await getter();
    console.log(`${value.toString()} ==? ${expected.toString()}`);
    if (value.toString() === expected.toString()) {
      return;
    }
  }
  throw new Error(`${value.toString()} never was ${expected.toString()}`);
}

describe("PeggyService", () => {
  let EthbridgeService: ReturnType<typeof createEthbridgeService>;

  beforeEach(async () => {
    EthbridgeService = createEthbridgeService({
      sifApiUrl: "http://localhost:1317",
      sifWsUrl: "ws://localhost:26667/nosocket",
      sifChainId: "sifchain",
      bridgebankContractAddress: config.bridgebankContractAddress,
      bridgetokenContractAddress: (EROWAN as Token).address,
      getWeb3Provider,
    });
  });

  // Values here are not working so skipping
  test.skip("lock and burn eth <-> ceth", async () => {
    // Setup services
    const sifService = await createTestSifService(akasha);
    const ethService = await createTestEthService();

    function getEthAddress() {
      return ethAccounts[0].public;
    }

    function getSifAddress() {
      return akasha.address;
    }

    async function getEthBalance() {
      const bals = await ethService.getBalance(getEthAddress(), ETH);
      return bals[0].toBaseUnits();
    }

    async function getCethBalance() {
      const bals = await sifService.getBalance(getSifAddress(), CETH);
      return bals[0].toBaseUnits();
    }

    const web3 = new Web3(await getWeb3Provider());

    // Check the balance
    const cethBalance = await getCethBalance();

    // Send funds to the smart contract
    await new Promise<void>(async done => {
      EthbridgeService.lockToSifchain(
        getSifAddress(),
        AssetAmount(ETH, "2"),
        10
      )
        .onTxHash(() => {
          advanceBlock(100);
        })
        .onComplete(() => {
          done();
        })
        .onError(err => {
          throw err.payload;
        });
    });

    const expectedCethAmount = JSBI.add(
      cethBalance,
      AssetAmount(ETH, "2").toBaseUnits()
    );

    await waitFor(
      async () => await getCethBalance(),
      expectedCethAmount,
      "expectedCethAmount"
    );

    const recipientBalanceBefore = await getEthBalance();

    const amountToSend = AssetAmount(CETH, "2");
    const feeAmount = AssetAmount(
      Asset.get("ceth"),
      JSBI.BigInt("16164980000000000")
    );

    const message = await EthbridgeService.burnToEthereum({
      fromAddress: getSifAddress(),
      assetAmount: amountToSend,
      feeAmount,
      ethereumRecipient: getEthAddress(),
    });

    // Message has the expected format
    const ethereumChainId = await web3.eth.net.getId();
    expect(message.value.msg).toEqual([
      {
        type: "ethbridge/MsgBurn",
        value: {
          cosmos_sender: getSifAddress(),
          amount: "2000000000000000000",
          symbol: "ceth",
          ethereum_chain_id: `${ethereumChainId}`,
          ethereum_receiver: getEthAddress(),
        },
      },
    ]);

    await sifService.signAndBroadcast(message.value.msg);

    const expectedEthAmount = JSBI.add(
      recipientBalanceBefore,
      JSBI.BigInt("2000000000000000000")
    ).toString();

    await waitFor(
      async () => await getEthBalance(),
      expectedEthAmount,
      "expectedEthAmount"
    );
  });

  test("rowan -> erowan -> rowan", async () => {
    const sifService = await createTestSifService(juniper);
    const ethService = await createTestEthService();
    const web3 = new Web3(await getWeb3Provider());

    function getEthAddress() {
      return ethAccounts[2].public;
    }

    function getSifAddress() {
      return juniper.address;
    }

    async function getERowanBalance() {
      const bals = await ethService.getBalance(getEthAddress(), EROWAN);
      return bals[0].toBaseUnits();
    }

    async function getRowanBalance() {
      const bals = await sifService.getBalance(getSifAddress(), ROWAN);
      return bals[0].toBaseUnits();
    }

    ////////////////////////
    // Rowan -> eRowan
    ////////////////////////

    // First get balance in ethereum
    const startingERowanBalance = await getERowanBalance();

    // lock Rowan to eRowan
    const sendRowanAmount = AssetAmount(ROWAN, "100");

    const msg = await EthbridgeService.lockToEthereum({
      fromAddress: getSifAddress(),
      assetAmount: sendRowanAmount,
      ethereumRecipient: getEthAddress(),
      feeAmount: AssetAmount(
        Asset.get("ceth"),
        JSBI.BigInt("18332015000000000")
      ),
    });

    const ethereumChainId = await web3.eth.net.getId();
    expect(msg.value.msg).toEqual([
      {
        type: "ethbridge/MsgLock",
        value: {
          amount: "100000000000000000000",
          ceth_amount: "18332015000000000",
          cosmos_sender: getSifAddress(),
          ethereum_chain_id: `${ethereumChainId}`,
          ethereum_receiver: getEthAddress(),
          symbol: "rowan",
        },
      },
    ]);

    await sifService.signAndBroadcast(msg.value.msg);

    await sleep(2000);

    const expectedERowanBalance = JSBI.add(
      startingERowanBalance,
      sendRowanAmount.toBaseUnits()
    );

    await waitFor(
      async () => await getERowanBalance(),
      expectedERowanBalance,
      "expectedERowanBalance"
    );

    ////////////////////////
    // eRowan -> Rowan
    ////////////////////////
    const startingRowanBalance = await getRowanBalance();

    const sendERowanAmount = AssetAmount(EROWAN, "10");

    // Burn eRowan to Rowan
    await new Promise<void>((done, reject) => {
      EthbridgeService.burnToSifchain(
        getSifAddress(),
        sendERowanAmount,
        50,
        getEthAddress()
      )
        .onTxHash(() => {
          advanceBlock(52);
        })
        .onComplete(async () => {
          console.log("COMPLETE!! 🍾");
          done();
        })
        .onError(err => {
          reject(err);
        });
    });

    await sleep(2000);

    // wait for the balance to change
    const expectedRowanBalance = JSBI.add(
      startingRowanBalance,
      sendERowanAmount.toBaseUnits()
    );

    await waitFor(
      async () => await getRowanBalance(),
      expectedRowanBalance,
      "expectedRowanBalance"
    );
  });
});
