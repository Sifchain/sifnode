# Peggy 2.0 PRD

---

**Note:**
This is documentation in progress.
It was created with [docsify](https://docsify.js.org/#/?id=docsify)
To use it:

```
cd docs/peggy
npm install
./node_modules/docsify-cli/bin/docsify serve ./docs
browser http://localhost:3000/
```

---

## Problem Statement

The Ethereum scaling problem is a well known issue for traders resulting in long confirmation times and exorbitant fees. The recent rise of non-Ethereum-based exchanges has shown that there is considerable demand for defi products on a more scalable platform. Sifchain, our application-specific blockchain, is an ever more efficient platform for value exchange. That said, the vast majority of development is still on Ethereum. To fully realize the advantages of all platforms, bridges are the keystone in creating a united defi ecosystem.

## One Liner

Peggy 2.0 is a system that facilitates non-permissioned token transfer between Sifchain and EVM chains built with extensibility to multiple chains.

## Goals

Incremental increases in the following

- TVL 
- Trading Volume 
- Weekly Active Addresses 
- Number of non-Ethereum tokens imported into Sifchain

## Use Cases

<table>
  <tbody>
    <tr>
      <td>Scenario</td>
      <td>User</td>
      <td>Goal</td>
      <td>Use Case</td>
    </tr>
    <tr>
      <td>1</td>
      <td>Poolers, Traders, Importer, Exporter</td>
      <td>Search for the right token to interact with</td>
      <td>
        <ol>
          <li>User searches for the token within approved list </li>
          <li>User has the option to examine exact token details to
            ensure they are using the right token</li>
          <li>User has the option to enter in a precise denom/contract address if
            token cannot be found on list</li>
        </ol>
        <a href="https://docs.google.com/document/d/1X8mnfAPVmK1_MfV7UJKjDRztoARVjhBt1CaqYAX0gvw/edit">Link to
          document</a>
      </td>
    </tr>
    <tr>
      <td>2</td>
      <td>Importer</td>
      <td>Move EVM Native Coins Into Sifchain (e.g. ETH-ETH, BSC-BNB)</td>
      <td>
        <ol>
          <li>User connects their Sifchain and EVM wallets </li>
          <li> User finds the correct token, ETH-ETH (Scenario 1) </li>
          <li> User
            imports ETH-ETH</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>3</td>
      <td>Importer</td>
      <td>Move EVM Wrapped Coins Into Sifchain (e.g. ETH-WBTC, BSC-ETH)</td>
      <td>
        <ol>
          <li>User connects their Sifchain and EVM wallets </li>
          <li>User finds the correct token, ETH-WBTC (Scenario 1) </li>
          <li> User
            imports ETH-WBTC</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>2</td>
      <td>Pooler</td>
      <td>Earn pooling rewards on EVM Native Coins</td>
      <td>
        <ol>
          <li>User connects their Sifchain and EVM wallets </li>
          <li> User finds the correct token, ETH-ETH (Scenario 1) </li>
          <li> User
            imports ETH-ETH </li>
          <li> User searches for the ETH-ETH pool (Scenario 1) </li>
          <li> User pools ETH-ETH</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>3</td>
      <td>Pooler</td>
      <td>Earn pooling rewards on EVM Wrapped Coins</td>
      <td>
        <ol>
          <li>User connects their BSC and EVM wallets </li>
          <li> User finds the correct token, ETH-WBTC (Scenario 1) </li>
          <li> User
            imports ETH-WBTC </li>
          <li> User searches for the ETH-WBTC pool (Scenario 1) </li>
          <li> User pools ETH-WBTC</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>4</td>
      <td>Trader</td>
      <td>Swap Ethereum Native for BSC Native (e.g. ETH-ETH for BSC-BNB)</td>
      <td>
        <ol>
          <li> User connects their BSC and Ethereum wallets </li>
          <li> User finds the correct tokens (Scenario 1) </li>
          <li> User swaps
            tokens</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>5</td>
      <td>Trader</td>
      <td>Swap Ethereum Native for BSC Wrapped (e.g. ETH-ETH for BSC-ETH)</td>
      <td>
        <ol>
          <li> User connects their BSC and Ethereum wallets </li>
          <li> User finds the correct tokens - can tell the difference
            between ETH-ETH and BSC-ETH (Scenario 1) </li>
          <li> User swaps tokens</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>6</td>
      <td>Exporter</td>
      <td>Export EVM Native to different EVM chain (e.g. ETH-ETH to BSC, BSC-BNB to Ethereum)</td>
      <td>
        <ol>
          <li> User moves EVM Native Coins Into Sifchain (Scenario 2) </li>
          <li> User finds the correct token, ETH-ETH (Scenario
            1) and selects destination network </li>
          <li> User conducts export</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>7</td>
      <td>Exporter</td>
      <td>Export EVM Wrapped to different EVM chain (e.g. ETH-WBTC to BSC, BSC-CAKE to Ethereum)</td>
      <td>
        <ol>
          <li> User moves EVM Wrapped Coins Into Sifchain (Scenario 2) </li>
          <li> User finds the correct token, ETH-WBTC (Scenario
            1) and selects destination network </li>
          <li> User conducts export</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>8</td>
      <td>Exporter</td>
      <td>Export Currently Supported Cosmos Native to EVM chain (e.g. IBC-AKASH to BSC)</td>
      <td>
        <ol>
          <li> User connects Sifchain and EVM wallets </li>
          <li> User finds the correct token (Scenario 1) and selects destination
            network </li>
          <li> User conducts export</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>9</td>
      <td>Delegator</td>
      <td>Delegate Rowan to a validator so that they can earn Rowan</td>
      <td>In Keplr Wallet: <ol>
          <li>User connects wallet and delegates via Keplr interface </li>
          <li> User can see claimable Rowan in
            wallet </li>
          <li> User claims Rowan in wallet</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>10</td>
      <td>Liquidity Provider</td>
      <td>Claim rewards for providing liquidity to Sifchain pools</td>
      <td>
        <ol>
          <li> User connects Keplr wallet </li>
          <li> On https://dex.sifchain.finance/#/rewards, user views rewards as proscribed
            by the current LM mining program </li>
          <li> User can claim reward</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>11</td>
      <td>Exporter</td>
      <td>Export EVM Native to Cosmos chain (e.g. ETH-ETH to Akash)</td>
      <td>
        <ol>
          <li> User connects Sifchain and Cosmos wallets - Double check on this w/Casey should be covered by IBC </li>
          <li> User
            finds the correct token (Scenario 1) and selects destination network </li>
          <li> User conducts export</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>12</td>
      <td>Exporter</td>
      <td>Export EVM Wrapped to Cosmos chain (e.g. ETH-WBTC to BSC, BSC-CAKE to Ethereum)</td>
      <td>
        <ol>
          <li> User moves EVM Wrapped Coins Into Sifchain (Scenario 2) </li>
          <li> User finds the correct token, ETH-WBTC (Scenario
            1) and selects destination network </li>
          <li> User conducts export</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>13</td>
      <td>Exporter</td>
      <td>Export Cosmos tokens to Cosmos chain (e.g. IBC-AKASH to Cosmos Hub)</td>
      <td>
        <ol>
          <li> User moves Cosmos Tokens into Sifchain </li>
          <li> User finds the correct token, IBC-AKASH (Scenario 1) and selects
            destination network </li>
          <li> User conducts export</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>14</td>
      <td>Exporter</td>
      <td>Export Sifchain Wrapped IBC token to Ethereum for the first time</td>
      <td>
        <ol>
          <li> User connects Keplr and Ethereum Wallets </li>
          <li> User finds the correct token (Scenario 1), and attempts to
            conduct export </li>
          <li> User is prompted to pay (in Sifchain ETH) for the full cost of exporting a new ERC20 token to
            Ethereum</li>
        </ol>
      </td>
    </tr>
    <tr>
      <td>15</td>
      <td>Engineer</td>
      <td>Pause relayer</td>
      <td>Pause-permissioned user should be able to pause the relayers, but not unpause them</td>
    </tr>
    <tr>
      <td>16</td>
      <td>Engineer</td>
      <td>Unpause relayer</td>
      <td>Unpause-permissioned user should be able to pause and unpause the relayers</td>
    </tr>
  </tbody>
</table>

## Peggy 2.0 Features

- Updated Sifchain Token Denoms - Extensibility 
- Stateless Signature Aggregation - Gas Savings
- Transaction Ordering - Stability
- EVM to EVM double pegging - Extensibility
- Removal of Ethereum whitelists - Usability
- Individually Batched Imports/Exports - Gas Savings


## Peggy 2.1

- New EVM chain integration 


## Roadmap

- Initial Development
- Integration Testing Group 0
- Documentation for Auditing 
- Integration Testing Group 1 
- Merge Develop to Future/Peggy2
- Halborn Audit 1 
- Integration Testing Group 2 
- Migration/Load Testing 
- Migration
- Halborn Audit 2 
- Deployment/Devops
- UI Updates
- Sifnode Updates
- Data Pipeline Infrastructure Changes
- Launch Peggy 2.0
