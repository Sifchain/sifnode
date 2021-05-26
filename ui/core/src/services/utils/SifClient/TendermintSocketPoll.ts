import { EventEmitter2 } from "eventemitter2";
import axios from "axios";
export type TendermintSocketPoll = ReturnType<typeof TendermintSocketPoll>;

// Partial BlockData we need
export type BlockData = {
  result: {
    block: {
      header: {
        height: string; // number in string
      };
      data: {
        txs: null | string[]; // Array of hashes
      };
    };
  };
};

async function fetchBlock(url: string): Promise<BlockData> {
  const res = await axios.get(url);
  return res.data;
}

export type ITendermintSocketPoll = ReturnType<typeof TendermintSocketPoll>;

export function TendermintSocketPoll({
  apiUrl,
  fetcher = fetchBlock,
  pollInterval = 5000,
}: {
  apiUrl: string;
  fetcher?: typeof fetchBlock;
  pollInterval?: number;
}) {
  const emitter = new EventEmitter2();

  async function pollBlock(height?: number) {
    const query = typeof height !== "undefined" ? `?height=${height}` : "";
    return await fetcher(`${apiUrl.replace(/\/$/, "")}/block${query}`);
  }

  // Process a block and emit events based on that block
  function processData(blockData: BlockData) {
    emitter.emit("NewBlock", blockData);

    const txs = blockData.result.block.data.txs;

    if (txs) {
      txs.forEach((tx) => {
        // TODO: Not sure if we should/can add more tx information here as all we have is the encoded tx - can we decode it? need to look into it
        emitter.emit("Tx", tx);
      });
    }

    // Return processed blockheight
    return parseInt(blockData.result.block.header.height);
  }

  function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  /**
   * Get all blocks that havent been processed since last known blockheight
   * @param height last processed blockheight or null for no new blockheight
   */
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

    // There are blocks to be processed build up a list of them
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
    // If already polling dont poll again
    if (polling) return;

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

// Make this a singleton to avoid multiple polling
let instance: ITendermintSocketPoll;
export function createTendermintSocketPoll(apiUrl: string) {
  if (!instance) {
    instance = TendermintSocketPoll({ apiUrl });
  }
  return instance;
}
