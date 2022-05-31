import { Client } from "pg";

// Goal:
/**
 * Goal:
 * 1. Create a function f(cosmosAddress) => {SifchainEvents}
 * Stretch:
 * 2. Use this function to iterate through all cosmosAddress seen
 * 3. Store this into a graph database
 */

type cosmosAddress = string;
type denom = string;

type raw_burn = {
  type: string;
  burner: cosmosAddress;
  // This ISNT denom. it is amount + Denom wtf .v.
  amount: denom;
};

type raw_coin_received = {
  receiver: cosmosAddress;
  // This ISNT denom. it is amount + Denom wtf .v.
  // 529846562ibc/C5C8682EB9AA1313EF1B12C991ADCDA465B80C05733BFB2972E2005E01BCE459
  amount: denom;
};

type raw_ibc_transfer = {
  sender: cosmosAddress;
  receiver: cosmosAddress;
};

// I Create this type so if we want to attach metadata we can
type IBCTransferEvent = raw_ibc_transfer & { amount: denom };

let pgClient: Client;

// Return type?
function fetchTransfer(cosmosAddress: cosmosAddress) {}

async function retrieveRawIBCEvent(): Promise<raw_ibc_transfer[]> {
  let query: string = "SELECT * from events_audit limit 1";

  console.log("Querying", pgClient);
  pgClient.query(
    query,
    (out, err) => {
      console.log(out);
      console.log(err);
    }
    // values: [],
  );
  // .then((output) => console.log("Output", output))
  // .catch((err) => console.error("Error rawibc query", err));

  return [];
}

// function rawIbcToIbc(
//   logs: (raw_ibc_transfer | raw_burn | raw_coin_received)[]
// ): ibc_transfer[] {
//   logs.filter(log => ( log instanceof raw_ibc_transfer ));
//   return null;
// }

// [burn | coin_received | coin_spent | ibc_transfer | message | send_packet | transfer ] => ibc_transfer

// function rawIBCTransferToInternalTransfer

// Return type?
// function fetchEVMExports(cosmosAddress: string) {}

function fetchIBCExports(cosmosAddress: cosmosAddress) {}

function fetchIBCImports(cosmosAddress: cosmosAddress) {}

function filterIBCByCosmosAddress() {}

function fetchSwapEvents(cosmosAddress: cosmosAddress) {}

function initPg(): Client {
  // const pgClient = new Client();
  pgClient = new Client();
  pgClient.connect();
  console.log("Connected to tsdb");
  return pgClient;
}

async function endPg(): Promise<void> {
  console.log("Closing connection");
  await pgClient.end();
}

async function main(): Promise<void> {
  try {
    initPg();
    retrieveRawIBCEvent();
  } catch (error) {
    console.log("Had err");
  } finally {
    endPg();
  }
}

main()
  .then(() => console.log("filterCosmosEvent.ts exiting"))
  .catch((err) => console.error("Error in filterCosmosEvent.ts", err));
