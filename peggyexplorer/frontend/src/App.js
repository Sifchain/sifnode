import React, { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';
import { TransactionDescription } from 'ethers/lib/utils';

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
  const createClaims = cosmosTxs.map(tx => {
    let log;
    try {
      log = JSON.parse(tx && tx.tx_result && tx.tx_result.log);
    } catch (e) {
    }
    return { tx, log };
  }).filter(({ log }) => {
    if (log !== undefined) {
      return log[0].events && log[0].events[0].type === 'create_claim';
    }
    return false;
  });
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
          Cosmos Events: {cosmosTxsChecked}
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
            </tr>
          </thead>
          <tbody>
            {locksOnly.map(e => {
              const ropstenURL = `https://ropsten.etherscan.io/address/${e.transactionHash}`;
              window.txtest = e;
              window.createClaims = createClaims;
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
              </tr>
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default App;
