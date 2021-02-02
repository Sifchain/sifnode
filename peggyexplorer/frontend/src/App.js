import React, { useState, useEffect } from 'react';
import './App.css';
const _ = require('lodash');

const SERVER_URL = "http://localhost:5000/dump";

function App() {

  const [dump, setDump] = useState({});

  const blocksChecked = dump && dump.allEthereumBlocksChecked && dump.allEthereumBlocksChecked.length;
  const ethereumEvents = dump && dump.allEthereumEvents || [];
  const ethereumEventsChecked = dump && dump.allEthereumEvents && dump.allEthereumEvents.length;
  const cosmosTxs = dump && dump.allCosmosTxs || [];
  const cosmosTxsChecked = dump && dump.allCosmosTxs && dump.allCosmosTxs.length;

  const loadStateDump = _ => {
    fetch(SERVER_URL)
      .then(res => res.json())
      .then(data => setDump(data));
  }

  useEffect(() => {
    // Update the document title using the browser API
    loadStateDump();
  }, []);

  const locksOnly = ethereumEvents.filter(e => e.event == "LogLock");
  const createClaimTxs = cosmosTxs.map(tx => {
    let fullLog;
    try {
      fullLog = JSON.parse(tx && tx.tx_result && tx.tx_result.log);
    } catch (e) {
    }
    return { tx, fullLog };
  }).map(({ tx, fullLog }) => {
    if (fullLog !== undefined) {
      const claimEvents = fullLog.reduce((accum, singleLog) => {
        accum.push(...singleLog.events.filter(event => event.type === 'create_claim'));
        return accum;
      }, []);
      return { tx, claimEvents };
    }
  });
  let createClaimEvents = createClaimTxs.reduce((accum, { claimEvents, tx }) => {
    const events = claimEvents.map(event => ({ event, tx }));
    accum.push(...events);
    return accum;
  }, []);
  createClaimEvents = createClaimEvents.map(claimEvent => {
    const { event, tx } = claimEvent;
    const { type, attributes } = event;
    const betterAttributes = attributes.reduce((accum, attribute) => {
      accum[attribute.key] = attribute.value;
      return accum;
    }, {});
    const betterEvent = Object.assign({}, { type, attributes: betterAttributes })
    return Object.assign({}, { event: betterEvent, tx });
  });
  const createClaimEventsByNonce = _.groupBy(createClaimEvents, claimEvent => claimEvent.event.attributes.nonce);
  return (
    <div className="App">
      <header className="App-header">
        <span>
          Blocks Checked: {blocksChecked}
        </span>
        <span>
          Ethereum Events: {ethereumEventsChecked}
        </span>
        <span>
          Cosmos Txs: {cosmosTxsChecked}
        </span>
        <button onClick={loadStateDump}>Refresh</button>
      </header>
      <div>
        Ethereum locks:
        <table>
          <thead>
            <tr>
              <th>
                etherum tx hash
              </th>
              <th>
                event
              </th>
              <th>
                bridge nonce
              </th>
              <th>
                msgs relayed to cosmos
              </th>
            </tr>
          </thead>
          <tbody>
            {locksOnly.map(e => {
              const ropstenURL = `https://ropsten.etherscan.io/address/${e.transactionHash}`;
              const claimEvents = createClaimEventsByNonce[e.returnValues._nonce];
              const eventDetails = claimEvents && claimEvents.map(claimEvent => {
                return { txHash: claimEvent.tx.hash, validator: claimEvent.event.attributes.validator_address }
              });
              window.createClaimEventsByNonce = createClaimEventsByNonce;
              return <tr key={e.returnValues._nonce}>
                <td>
                  <a href={ropstenURL}>{e.transactionHash.slice(0, 10)}</a>
                </td>
                <td>
                  {e.event}
                </td>
                <td>
                  {e.returnValues._nonce}
                </td>
                <td>
                  {eventDetails.map(detail => <div key={detail.validator}>
                    txHash: {detail.txHash} <br />
                    validator: {detail.validator}
                  </div>)}
                </td>
              </tr>
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default App;
