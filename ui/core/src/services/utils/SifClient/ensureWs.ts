import ReconnectingWebSocket from "reconnecting-websocket";

// Pool socket connections
const wsPool: { [s: string]: ReconnectingWebSocket } = {};

export function ensureWs(wsUrl: string) {
  if (!wsPool[wsUrl]) {
    wsPool[wsUrl] = new ReconnectingWebSocket(wsUrl);
  }

  return wsPool[wsUrl];
}
