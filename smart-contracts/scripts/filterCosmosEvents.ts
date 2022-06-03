import { Pool } from "pg";

// Goal:
/**
 * Goal:
 * 1. Create a function f(cosmosAddress) => {SifchainEvents}
 * Stretch:
 * 2. Use this function to iterate through all cosmosAddress seen
 * 3. Store this into a graph database
 */
const TABLE_NAME = "events_audit";

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

interface raw_ibc_transfer {
  sender: cosmosAddress;
  receiver: cosmosAddress;
}

interface raw_transfer {
  receipient: cosmosAddress;
  sender: cosmosAddress;
  amount: string;
}

interface raw_send_packet {
  packet_data: {
    amount: number;
    denom: string;
  };
}

// I Create this type so if we want to attach metadata we can
// TODO: Add Time
type IBCTransferEvent = {
  ibc_transfer: raw_ibc_transfer;
  amount: number;
  denom: string;
  height: number;
};

type postgresLog = {
  type: string;
  attributes: {
    key: string;
    value: string;
  }[];
};

function etlrawIbc(ibc_transfer_entry: postgresLog): raw_ibc_transfer {
  if (ibc_transfer_entry.type != "ibc_transfer") {
    throw new Error("Invalid type");
  }

  let sender: cosmosAddress = "";
  let receiver: cosmosAddress = "";

  for (let attribute of ibc_transfer_entry.attributes) {
    if (attribute.key == "sender") {
      sender = attribute.value;
    } else if (attribute.key == "receiver") {
      receiver = attribute.value;
    }
  }

  if (sender == "" || receiver == "") {
    throw new Error("Sender/Receiver is empty");
  }

  return {
    sender: sender,
    receiver: receiver,
  };
}

function etlRawPacketData(send_packet_entry: postgresLog): raw_send_packet {
  if (send_packet_entry.type != "send_packet") {
    throw new Error("Invalid message type, expected send_packet");
  }

  for (const attribute of send_packet_entry.attributes) {
    if (attribute.key == "packet_data") {
      const entry = JSON.parse(attribute.value);
      return {
        packet_data: {
          amount: entry.amount,
          denom: entry.denom,
        },
      };
    }
  }
  throw new Error(
    "Expected [packet_data] to be an attribute, did not encounter"
  );
}

/**
 *
 * @param row Postgres Row. An object containing columns as object property
 * @returns
 */
function rowToIBCTransferEvent(row: any): IBCTransferEvent {
  // console.log("Received row:", JSON.stringify(row));
  let ibc_transfer: raw_ibc_transfer;
  let packet_data: raw_send_packet;

  const logsArray = row.log;

  for (const log of logsArray) {
    if (log.type == "ibc_transfer") {
      ibc_transfer = etlrawIbc(log);
    }
    if (log.type == "send_packet") {
      packet_data = etlRawPacketData(log);
    }
  }

  const ibcTransferEvent: IBCTransferEvent = {
    ibc_transfer: ibc_transfer!,
    amount: packet_data!.packet_data.amount,
    denom: packet_data!.packet_data.denom,
    height: row.height,
  };

  return ibcTransferEvent;
}

async function retrieveRawIBCEvent(pool: Pool): Promise<IBCTransferEvent[]> {
  const output: IBCTransferEvent[] = [];
  const TIME_RANGE_IN_DAYS = 30;

  const query: string =
    "SELECT time, type, height, log from " +
    TABLE_NAME +
    " where type = 'ibc_transfer' and " +
    " time > NOW() - interval '" +
    TIME_RANGE_IN_DAYS +
    " days' " +
    "";
  // " limit 5";

  // console.log("Querying", pool);
  const postgresOutput = await pool.query({
    text: query,
  });

  // return postgresOutput.rows.map((row) => row.log).map(rowToIBCTransferEvent);
  return postgresOutput.rows.map(rowToIBCTransferEvent);

  for (const row of postgresOutput.rows) {
    // console.log(row);
    const ibcTransferEvent: IBCTransferEvent = rowToIBCTransferEvent(row.log);
    output.push(ibcTransferEvent);
  }

  // console.log(output);
  return output;
}

function fetchIBCExports(
  cosmosAddress: cosmosAddress,
  ibcTransfers: IBCTransferEvent[]
) {
  const relevant = ibcTransfers.filter(
    (transfer) => transfer.ibc_transfer.sender == cosmosAddress
  );
  console.log("IBCExport:", relevant);
}

function fetchIBCImports(
  cosmosAddress: cosmosAddress,
  ibcTransfers: IBCTransferEvent[]
) {
  const relevant = ibcTransfers.filter(
    (transfer) => transfer.ibc_transfer.receiver == cosmosAddress
  );
  console.log("IBCImports:", relevant);
}

function filterIBCByCosmosAddress() {}

async function initPg(): Promise<Pool> {
  const pool = new Pool({
    keepAliveInitialDelayMillis: 180_000,
    idle_in_transaction_session_timeout: 180_000,
    connectionTimeoutMillis: 180_000,
    idleTimeoutMillis: 180_000,
    ssl: false,
  });
  await pool.connect();
  console.log("Connected to tsdb");
  return pool;
}

async function findTransfers(pool: Pool, addresses: string[]) {
  const query =
    "select bn_recipient, bn_sender, bn_amount, bn_token, log from events_audit where type in ('burn', 'lock') and log -> 0 -> 'attributes' -> 2 ->> 'key' = 'ethereum_chain_id' and time > NOW() - interval '30 days';";
  const values = [];
  const result = await pool.query(query);
  return result;
}

async function lockAndBurn(pool: Pool) {
  const query =
    "select bn_recipient, bn_sender, bn_amount, bn_token, log from events_audit where type in ('burn', 'lock') and log -> 0 -> 'attributes' -> 2 ->> 'key' = 'ethereum_chain_id' and time > NOW() - interval '30 days';";
  const out = await pool.query(query);
  console.log(out);
}

async function endPg(pool: Pool): Promise<void> {
  console.log("Closing connection");
  await pool.end();
}

async function main(): Promise<void> {
  let pool: Pool;

  const addressOfInterest: cosmosAddress[] = [
    "sif1vmzzvtwr5dl2dumh0l65fsvfdf6j2etv2tt4xy",
    "sif164lcv2wxzyzy6g4ea7c2jrwjhqukfll8jr6wa2",
  ];

  try {
    pool = await initPg();
  } catch (error) {
    console.error("Error connecting occurred: ", error);
    return;
  }

  try {
    const ibcTransferEvents = await retrieveRawIBCEvent(pool);
    fetchIBCExports(addressOfInterest[0], ibcTransferEvents);
    fetchIBCImports(addressOfInterest[0], ibcTransferEvents);
  } catch (error) {
    console.log("Had error: ", error);
  } finally {
    endPg(pool);
  }
}

main()
  .then(() => console.log("filterCosmosEvent.ts exiting"))
  .catch((err) => console.error("Error in filterCosmosEvent.ts", err));
