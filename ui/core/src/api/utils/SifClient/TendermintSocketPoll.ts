import { EventEmitter2 } from "eventemitter2";

export type TendermintSocketPoll = ReturnType<
  typeof TendermintSocketPoll
>;

export function TendermintSocketPoll({ apiUrl }: { apiUrl: string }) {
  const emitter = new EventEmitter2();

  const fetchNewBlocks = async () => {
    const blockResponse = await fetch(`${apiUrl}blocks/latest`);
    const data = await blockResponse.json();
    const eventType = "NewBlock";
    emitter.emit(eventType, data);
  }
  return {
    on(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void) {
      if (!emitter.hasListeners(event)) {
        console.log('fetching url');
        this.fetchNewBlocksInterval = setInterval(fetchNewBlocks, 5000)
      }
      emitter.on(event, handler);
    },
    off(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void) {
      emitter.off(event, handler);
      clearInterval(this.fetchNewBlocksInterval)
    },
  };
}

export function createTendermintSocketPoll(apiUrl: string) {
  return TendermintSocketPoll({ apiUrl });
}
