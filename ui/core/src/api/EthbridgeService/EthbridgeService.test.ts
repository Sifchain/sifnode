import createEthbridgeService from ".";
import { Asset, AssetAmount } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";

import { advanceBlock } from "../../test/utils/advanceBlock";
import { createWaitForBalance } from "../../test/utils/waitForBalance";
import { akasha } from "../../test/utils/accounts";

import { createTestSifService } from "../../test/utils/services";
import { getBalance, getTestingToken } from "../../test/utils/getTestingToken";
import Web3 from "web3";

describe("PeggyService", () => {
  let EthbridgeService: ReturnType<typeof createEthbridgeService>;

  let ETH: Asset;
  let CETH: Asset;

  beforeEach(async () => {
    ETH = getTestingToken("ETH");
    CETH = getTestingToken("CETH");

    EthbridgeService = createEthbridgeService({
      sifApiUrl: "http://localhost:1317",
      sifWsUrl: "ws://localhost:26667/nosocket",
      sifChainId: "sifchain",
      bridgebankContractAddress: "0xf204a4Ef082f5c04bB89F7D5E6568B796096735a",
      getWeb3Provider,
    });
  });

  async function testLockEth() {
    // get sif balance
    const sifService = createTestSifService(akasha);

    const waitForBalance = createWaitForBalance(sifService);
    await waitForBalance("ceth", "1000000000", akasha.address);

    await new Promise<void>(async (done) => {
      EthbridgeService.lock(akasha.address, AssetAmount(ETH, "2"), 10)
        .onTxEvent((evt) => {
          console.log(evt);
        })
        .onComplete(async () => {
          // Not testing balances because we have no
          // way to correlate against transaction
          await waitForBalance("ceth", "2000000001000000000", akasha.address);
          done();
        })
        .onError((err) => {
          throw err.payload;
        });
      advanceBlock(200);
    });
  }

  test("lock eth", testLockEth);

  test.only("burn ceth", async () => {
    // we need to lock eth in the contract to test burn
    await testLockEth();

    const web3 = new Web3(await getWeb3Provider());
    const sifService = createTestSifService(akasha);
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
      assetAmount: AssetAmount(CETH, "2000000000000000000"),
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
});
