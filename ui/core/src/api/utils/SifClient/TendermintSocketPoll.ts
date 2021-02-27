import { EventEmitter2 } from "eventemitter2";
import axios from "axios";
export type TendermintSocketPoll = ReturnType<typeof TendermintSocketPoll>;

type BlockData = {
  result: {
    block: {
      header: {
        height: string; // number in string
      };
      txs: null | string[]; // Array of hashes
    };
  };
};

async function fetchBlock<T extends BlockData>(url: string): Promise<T> {
  const res = await axios.get(url);
  return res.data;
}

export function TendermintSocketPoll({
  apiUrl,
  fetcher = fetchBlock,
}: {
  apiUrl: string;
  fetcher?: typeof fetchBlock;
}) {
  const pollInterval = 5000;
  const emitter = new EventEmitter2();

  async function pollBlock(height?: number) {
    const query = typeof height !== "undefined" ? `?height=${height}` : "";
    return await fetcher(`${apiUrl}block${query}`);
  }

  // Process a block and emit events based on that block
  function processData(blockData: BlockData) {
    emitter.emit("NewBlock", blockData);

    const txs = blockData.result.block.txs;
    if (txs) {
      txs.forEach(tx => {
        // TODO: Not sure if we should/can add more tx information here as all we have is the encoded tx - can we decode it? need to look into it
        emitter.emit("Tx", tx);
      });
    }

    // Return processed blockheight
    return parseInt(blockData.result.block.header.height);
  }

  function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async function getBlocksToProcess(height: number | null) {
    const blockData = await pollBlock();
    const newHeight = parseInt(blockData.result.block.header.height);

    // If height is null this is the first poll so process this block
    if (height === null) {
      return [blockData];
    }

    // no new data don't process any blocks
    if (newHeight === height) {
      return [];
    }

    // Build up a list of blocks to process in order
    let heightDiff = newHeight - height;
    const blocks = [blockData];
    for (let i = heightDiff - 1; i > 0; --i) {
      const interimData = await pollBlock(height + i);
      blocks.unshift(interimData);
    }

    return blocks;
  }

  let polling = false;
  async function startPoll() {
    polling = true;
    let height: number | null = null;

    // Loop while we are polling
    while (polling) {
      // First we get a list of blocks to process
      const blocks = await getBlocksToProcess(height);

      // Then we process them updating the height
      for (let block of blocks) {
        height = processData(block);
      }

      // Then let's wait for a poll interval
      await sleep(pollInterval);
    }
  }
  function stopPoll() {
    polling = false;
  }

  return {
    on(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void) {
      if (!emitter.hasListeners()) {
        startPoll();
      }
      emitter.on(event, handler);
    },
    off(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void) {
      emitter.off(event, handler);
      if (!emitter.hasListeners()) {
        stopPoll();
      }
    },
  };
}

export function createTendermintSocketPoll(apiUrl: string) {
  return TendermintSocketPoll({ apiUrl });
}
