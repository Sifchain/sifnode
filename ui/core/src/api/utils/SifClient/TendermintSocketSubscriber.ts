import ReconnectingWebSocket from "reconnecting-websocket";
import { EventEmitter2 } from "eventemitter2";
import { uniqueId } from "lodash";
import { ensureWs } from "./ensureWs";

// Helper to allow us to add listeners to the open websocket
// In kind of a synchronous looking way
function openWebsocket(ws: ReconnectingWebSocket) {
  const wsPromise = new Promise<ReconnectingWebSocket>((resolve) => {
    if (ws.readyState === ReconnectingWebSocket.OPEN) {
      resolve(ws);
      return;
    }

    ws.addEventListener("open", () => {
      resolve(ws);
    });
  });

  return wsPromise.then.bind(wsPromise);
}

export type TendermintSocketSubscriber = ReturnType<
  typeof TendermintSocketSubscriber
>;

// Simplify subscribing to Tendermintsocket
export function TendermintSocketSubscriber({ wsUrl }: { wsUrl: string }) {
  const emitter = new EventEmitter2();

  // This let's us wait until the websocket is open before subscribing to messages on it
  const _ws = ensureWs(wsUrl);
  const withWebsocket = openWebsocket(_ws);

  withWebsocket((ws) => {
    ws.addEventListener("message", (message) => {
      const data = JSON.parse(message.data);

      const eventData = data.result?.data;

      if (!eventData) return;
      // Get last part of Tendermint Tx eg. 'tendermint/event/Tx'
      const [eventType] = eventData.type.split("/").slice(-1);

      // console.log("Message received");
      // console.log({ eventType, eventData });

      emitter.emit(eventType, eventData);
    });
  });

  return {
    on(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void) {
      // If for error listen immediately
      if (event === "error") {
        _ws.addEventListener("error", handler);
        return;
      }

      if (!emitter.hasListeners(event)) {
        withWebsocket((ws) => {
          ws.send(
            JSON.stringify({
              jsonrpc: "2.0",
              method: "subscribe",
              id: uniqueId(),
              params: {
                query: `tm.event='${event}'`,
              },
            })
          );
        });
      }
      emitter.on(event, handler);
    },
    off(event: "Tx" | "NewBlock" | "error", handler: (event: any) => void) {
      emitter.off(event, handler);
    },
  };
}

export function createTendermintSocketSubscriber(wsUrl: string) {
  return TendermintSocketSubscriber({ wsUrl });
}
