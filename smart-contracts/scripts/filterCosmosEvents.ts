import { Client, Pool } from "pg";

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

// Return type?
function fetchTransfer(cosmosAddress: cosmosAddress) {}

async function retrieveRawIBCEvent(pool: Pool): Promise<raw_ibc_transfer[]> {
  let query: string = "SELECT * from events_audit limit 1";

  console.log("Querying", pool);
  const output = await pool.query(query);
  console.log(output);
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

async function initPg(): Promise<Pool> {
  const pool = new Pool({
    keepAliveInitialDelayMillis: 180_000,
    idle_in_transaction_session_timeout: 180_000,
    connectionTimeoutMillis: 180_000,
    idleTimeoutMillis: 180_000,
    ssl: false
  });
  await pool.connect();
  console.log("Connected to tsdb");
  return pool;
}

async function findTransfers(pool: Pool, addresses: string[]) {
  const query = "select bn_recipient, bn_sender, bn_amount, bn_token, log from events_audit where type in ('burn', 'lock') and log -> 0 -> 'attributes' -> 2 ->> 'key' = 'ethereum_chain_id' and time > NOW() - interval '30 days';";
  const values = [];
  const result = await pool.query(query);
  return result;
}

async function lockAndBurn(pool: Pool) {
  const query = "select bn_recipient, bn_sender, bn_amount, bn_token, log from events_audit where type in ('burn', 'lock') and log -> 0 -> 'attributes' -> 2 ->> 'key' = 'ethereum_chain_id' and time > NOW() - interval '30 days';";
  const out = await pool.query(query);
  console.log(out);
}

async function endPg(pool: Pool): Promise<void> {
  console.log("Closing connection");
  await pool.end();
}

async function main(): Promise<void> {
  let pool : Pool;
  try {
    pool = await initPg();
  } catch (error) {
    console.error("Error connecting occurred: ", error);
    return;
  } try {
    await retrieveRawIBCEvent(pool);
    await lockAndBurn(pool);
  } catch (error) {
    console.log("Had error: ", error);
} finally {
    endPg(pool);
  }
}

main()
  .then(() => console.log("filterCosmosEvent.ts exiting"))
  .catch((err) => console.error("Error in filterCosmosEvent.ts", err));
