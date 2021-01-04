import ReconnectingWebSocket from "reconnecting-websocket";
import { EventEmitter2 } from "eventemitter2";
import { uniqueId } from "lodash";
type Handler<T = any> = (a: T) => void;

export type TendermintSocketSubscriber = ReturnType<
  typeof TendermintSocketSubscriber
>;

// Helper to allow us to add listeners to the open websocket
// In kind of a synchronous looking way
function openWebsocket(ws: ReconnectingWebSocket) {
  const wsPromise = new Promise<ReconnectingWebSocket>((resolve) => {
    ws.addEventListener("open", () => {
      resolve(ws);
    });
  });

  return wsPromise.then.bind(wsPromise);
}

export function TendermintSocketSubscriber({ wsUrl }: { wsUrl: string }) {
  const emitter = new EventEmitter2();

  // This let's us wait until the websocket is open before subscribing to messages on it
  const _ws = new ReconnectingWebSocket(wsUrl);
  const withWebsocket = openWebsocket(_ws);

  withWebsocket((ws) => {
    ws.addEventListener("message", (message) => {
      console.log({ message });
      const data = JSON.parse(message.data);

      const eventData = data.result?.data;

      if (!eventData) return;
      // Get last part of Tendermint Tx eg. 'tendermint/event/Tx'
      const [eventType] = eventData.type.split("/").slice(-1);

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
  };
}

export function createTendermintSocketSubscriber(wsUrl: string) {
  return TendermintSocketSubscriber({ wsUrl });
}

/*

 {
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "query": "tm.event='Tx'",
    "data": {
      "type": "tendermint/event/Tx",
      "value": {
        "TxResult": {
          "height": "1815",
          "index": 0,
          "tx": "tAEoKBapCjBUNypPChSBOsIqYUkWVeJ4auDRJ9aZC/v73hIGCgRjYXRrGgYKBGNidGsiBDEwMDASEAoKCgVyb3dhbhIBMBDAmgwaagom61rphyEC9QSwUduyvjSdNKZaHsJZhFkcbB/hylEu0mVpE7hUCioSQIyI9VUWx+rR5K8mFcoMoQeO2YMmLKAUIBtHEEiG9GY8LFOpln0IPrd/gx1MgGSq7t4oRAHcoAyp7QBD4zIUwBg=",
          "result": {
            "log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"swap\"},{\"key\":\"sender\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"},{\"key\":\"sender\",\"value\":\"sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85\"},{\"key\":\"module\",\"value\":\"clp\"},{\"key\":\"sender\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"}]},{\"type\":\"swap\",\"attributes\":[{\"key\":\"swap_amount\",\"value\":\"992\"},{\"key\":\"liquidity_fee\",\"value\":\"0\"},{\"key\":\"trade_slip\",\"value\":\"0\"},{\"key\":\"height\",\"value\":\"1815\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85\"},{\"key\":\"sender\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"},{\"key\":\"amount\",\"value\":\"1000catk\"},{\"key\":\"recipient\",\"value\":\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\"},{\"key\":\"sender\",\"value\":\"sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85\"},{\"key\":\"amount\",\"value\":\"996cbtk\"}]}]}]",
            "gas_wanted": "200000",
            "gas_used": "83990",
            "events": [
              {
                "type": "message",
                "attributes": [
                  {
                    "key": "YWN0aW9u",
                    "value": "c3dhcA=="
                  }
                ]
              },
              {
                "type": "transfer",
                "attributes": [
                  {
                    "key": "cmVjaXBpZW50",
                    "value": "c2lmMXBqbTIyOHJzZ3dxZjIzYXJreDdsbTl5cGt5bWE3bXpyM3kybjg1"
                  },
                  {
                    "key": "c2VuZGVy",
                    "value": "c2lmMXN5YXZ5Mm5wZnl0OXRjbmNkdHNkemY3a255OWxoNzc3eXFjMm5k"
                  },
                  {
                    "key": "YW1vdW50",
                    "value": "MTAwMGNhdGs="
                  }
                ]
              },
              {
                "type": "message",
                "attributes": [
                  {
                    "key": "c2VuZGVy",
                    "value": "c2lmMXN5YXZ5Mm5wZnl0OXRjbmNkdHNkemY3a255OWxoNzc3eXFjMm5k"
                  }
                ]
              },
              {
                "type": "transfer",
                "attributes": [
                  {
                    "key": "cmVjaXBpZW50",
                    "value": "c2lmMXN5YXZ5Mm5wZnl0OXRjbmNkdHNkemY3a255OWxoNzc3eXFjMm5k"
                  },
                  {
                    "key": "c2VuZGVy",
                    "value": "c2lmMXBqbTIyOHJzZ3dxZjIzYXJreDdsbTl5cGt5bWE3bXpyM3kybjg1"
                  },
                  {
                    "key": "YW1vdW50",
                    "value": "OTk2Y2J0aw=="
                  }
                ]
              },
              {
                "type": "message",
                "attributes": [
                  {
                    "key": "c2VuZGVy",
                    "value": "c2lmMXBqbTIyOHJzZ3dxZjIzYXJreDdsbTl5cGt5bWE3bXpyM3kybjg1"
                  }
                ]
              },
              {
                "type": "swap",
                "attributes": [
                  {
                    "key": "c3dhcF9hbW91bnQ=",
                    "value": "OTky"
                  },
                  {
                    "key": "bGlxdWlkaXR5X2ZlZQ==",
                    "value": "MA=="
                  },
                  {
                    "key": "dHJhZGVfc2xpcA==",
                    "value": "MA=="
                  },
                  {
                    "key": "aGVpZ2h0",
                    "value": "MTgxNQ=="
                  }
                ]
              },
              {
                "type": "message",
                "attributes": [
                  {
                    "key": "bW9kdWxl",
                    "value": "Y2xw"
                  },
                  {
                    "key": "c2VuZGVy",
                    "value": "c2lmMXN5YXZ5Mm5wZnl0OXRjbmNkdHNkemY3a255OWxoNzc3eXFjMm5k"
                  }
                ]
              }
            ]
          }
        }
      }
    },
    "events": {
      "swap.swap_amount": [
        "992"
      ],
      "message.module": [
        "clp"
      ],
      "tx.hash": [
        "2DBBA4107D509167DE4F6EFA1D2CBC2681ACD6FCA5A0988CE2E7143F52558E14"
      ],
      "tx.height": [
        "1815"
      ],
      "message.sender": [
        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
        "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
      ],
      "swap.liquidity_fee": [
        "0"
      ],
      "swap.trade_slip": [
        "0"
      ],
      "swap.height": [
        "1815"
      ],
      "message.action": [
        "swap"
      ],
      "transfer.recipient": [
        "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
      ],
      "transfer.sender": [
        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
        "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"
      ],
      "transfer.amount": [
        "1000catk",
        "996cbtk"
      ],
      "tm.event": [
        "Tx"
      ]
    }
  }
}
 
 */
