import createEthbridgeService from ".";
import { Asset, AssetAmount } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";

import { advanceBlock } from "../../test/utils/advanceBlock";
import { createWaitForBalance } from "../../test/utils/waitForBalance";
import { akasha } from "../../test/utils/accounts";
import { createTestSifService } from "../../test/utils/services";
import {
  getBalance,
  getTestingToken,
  getTestingTokens,
} from "../../test/utils/getTestingToken";
import config from "../../config.localnet.json";
import Web3 from "web3";
import JSBI from "jsbi";

const [ETH, CETH, ATK, CATK] = getTestingTokens(["ETH", "CETH", "ATK", "CATK"]);

describe("PeggyService", () => {
  let EthbridgeService: ReturnType<typeof createEthbridgeService>;

  beforeEach(async () => {
    EthbridgeService = createEthbridgeService({
      sifApiUrl: "http://localhost:1317",
      sifWsUrl: "ws://localhost:26667/nosocket",
      sifChainId: "sifchain",
      bridgebankContractAddress: config.bridgebankContractAddress,
      getWeb3Provider,
    });
  });

  test("lock and burn eth <-> ceth", async () => {
    // Setup services
    const sifService = createTestSifService(akasha);
    const waitForBalance = createWaitForBalance(sifService);
    const web3 = new Web3(await getWeb3Provider());

    // Check the balance
    await waitForBalance(
      "ceth",
      "99999991700000000000000000000",
      akasha.address
    );

    // Send funds to the smart contract
    await new Promise<void>(async done => {
      EthbridgeService.lock(akasha.address, AssetAmount(ETH, "2"), 10)
        .onComplete(async () => {
          // Check they arrived
          await waitForBalance(
            "ceth",
            "99999991702000000000000000000",
            akasha.address
          );
          done();
        })
        .onError(err => {
          throw err.payload;
        });
      advanceBlock(100);
    });

    const accounts = await web3.eth.getAccounts();
    const ethereumRecipient = accounts[0];
    const senderBalanceBefore = getBalance(
      await sifService.getBalance(akasha.address),
      "ceth"
    ).amount.toString();

    const recipientBalanceBefore = await web3.eth.getBalance(ethereumRecipient);

    console.log({
      ethereumRecipient,
      recipientBalanceBefore,
      senderBalanceBefore,
    });

    const ethereumChainId = await web3.eth.net.getId();
    const message = await EthbridgeService.burn({
      fromAddress: akasha.address,
      assetAmount: AssetAmount(CETH, "2"),
      ethereumRecipient,
    });

    // Message has the expected format
    expect(message).toEqual({
      type: "cosmos-sdk/StdTx",
      value: {
        msg: [
          {
            type: "ethbridge/MsgBurn",
            value: {
              cosmos_sender: akasha.address,
              amount: "2000000000000000000",
              symbol: "ceth",
              ethereum_chain_id: `${ethereumChainId}`,
              ethereum_receiver: ethereumRecipient,
            },
          },
        ],
        fee: {
          amount: [],
          gas: "200000",
        },
        signatures: null,
        memo: "",
      },
    });

    console.log(
      "Message to be broadcast:\n\n",
      JSON.stringify(message.value.msg, null, 2)
    );

    await sifService.signAndBroadcast(message.value.msg);

    const recipientBalanceAfter = await web3.eth.getBalance(ethereumRecipient);

    console.log({
      ethereumRecipient,
      recipientBalanceBefore,
      senderBalanceBefore,
      recipientBalanceAfter,
    });
  });

  test.todo("lock and burn atk <-> catk");
});
