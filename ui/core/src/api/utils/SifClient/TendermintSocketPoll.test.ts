import { sleep } from "../../../test/utils/sleep";
import { TendermintSocketPoll, BlockData } from "./TendermintSocketPoll";

// Simulate a fake block result from the server
function fakeBlockResult(
  height: number,
  txs: string[] | null = null,
): BlockData {
  return {
    result: {
      block: {
        header: {
          height: `${height}`,
        },
        data: {
          txs,
        },
      },
    },
  };
}

let latestBlock = 1;

// Simulate a fake block result from the server
const fetcher = async (requestUrl: string) => {
  // extract height from URL and create fake result
  const match = requestUrl.match(/\?height=(\d+)/);
  const num = (match && parseInt(match[1])) || latestBlock;

  // Have a couple of blocks contain txs
  // Should emit in order
  let tx = null;
  if (num === 4) {
    tx = ["baz", "bing", "boop"];
  }
  if (num === 7) {
    tx = ["foo", "bar"];
  }
  return fakeBlockResult(num, tx);
};

const newBlockHandler = jest.fn();
const txHandler = jest.fn();

beforeEach(() => {
  jest.resetAllMocks();
});

test("Fast interval has same results as slow interval", async () => {
  latestBlock = 1;
  const subscriber = TendermintSocketPoll({
    apiUrl: "http://localhost",
    fetcher,
    pollInterval: 10, // fast interval
  });

  subscriber.on("Tx", txHandler);
  subscriber.on("NewBlock", newBlockHandler);

  // set latest block to 10
  latestBlock = 10;

  // wait
  await sleep(1000);

  // Should be called exactly 10 times
  expect(newBlockHandler).toHaveBeenCalledTimes(10);
  expect(txHandler.mock.calls).toEqual([
    ["baz"],
    ["bing"],
    ["boop"],
    ["foo"],
    ["bar"],
  ]);
});

test("Slow interval has same results as fast interval", async () => {
  latestBlock = 1;
  const subscriber = TendermintSocketPoll({
    apiUrl: "http://localhost",
    fetcher,
    pollInterval: 1000, // fast interval
  });

  subscriber.on("Tx", txHandler);
  subscriber.on("NewBlock", newBlockHandler);

  // set latest block to 10
  latestBlock = 10;

  // wait
  await sleep(2000);

  // Should be called exactly 10 times
  expect(newBlockHandler).toHaveBeenCalledTimes(10);
  expect(txHandler.mock.calls).toEqual([
    ["baz"],
    ["bing"],
    ["boop"],
    ["foo"],
    ["bar"],
  ]);
});
