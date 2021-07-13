# REST API for Sifnode




### Base URL

  ```
  https://api.sifchain.finance
  ```




### Data types

  ```
  int64 Airdrop = 1
  int64 LiquidityMining = 2
  int64 ValidatorSubsidy = 3
  int64 DistributionType = Airdrop | LiquidityMining | ValidatorSubsidy

  CreateClaimReq: {
    BaseReq: BaseReq
    Signer: string
    ClaimType: DistributionType  
  }
  ```

  ```
  BaseReq: {
    from: string;
    chain_id: string;
    account_number?: string;
    sequence?: string;
  }
  ```

  ```
  RawPool: {
    external_asset: {
      source_chain: string;
      symbol: string;
      ticker: string;
    };
    native_asset_balance: string;
    external_asset_balance: string;
    pool_units: string;
    pool_address: string;
  }
  ```

  ```
  AminoMsg: {
    type: string;
    value: any;
  }
  ```

  ```
  BurnOrLockReq: {
    base_req: BaseReq;
    ethereum_chain_id: string;
    token_contract_address: string;
    cosmos_sender: string;
    ethereum_receiver: string;
    amount: string;
    symbol: string;
    ceth_amount: string;
  }
  ```
  
  ```
  LiquidityParams: {
    base_req: {
      from: string;
      chain_id: string;
    };
    external_asset: {
      source_chain: string;
      symbol: string;
      ticker: string;
    };
    native_asset_amount: string;
    external_asset_amount: string;
    signer: string;
  }
  ```

  ```
  SwapParams: {
    sent_asset: {
      symbol: string;
      ticker: string;
      source_chain: string;
    };
    received_asset: {
      symbol: string;
      ticker: string;
      source_chain: string;
    };
    base_req: {
      from: string;
      chain_id: string;
    };
    signer: string;
    sent_amount: string;
    min_receiving_amount: string;
  }
  ```

  ```
  LiquidityDetailsResponse: {
    result: {
      external_asset_balance: string;
      native_asset_balance: string;
      LiquidityProvider: {
        liquidity_provider_units: string;
        liquidity_provider_address: string;
        asset: {
          symbol: string;
          ticker: string;
          source_chain: string;
        };
      };
    };
    height: string;
  }
  ```

  ```
  RemoveLiquidityParams: {
    base_req: {
      from: string;
      chain_id: string;
    };
    external_asset: {
      source_chain: string;
      symbol: string;
      ticker: string;
    };
    w_basis_points: string;
    asymmetry: string;
    signer: string;
  }
  ```




<br />

# Dispensation module




<br />

### Create claim

* **URL**

  ```/dispensation/createClaim```

* **Method:**

  `POST`

* **Data Params**

  ```
  CreateClaimReq: { ... }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**

      ```
      { "type", "value" }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Create one of the predefined claim types.




<br />

# Ethbridge module




<br />

### Burn tokens

* **URL**

  ```/ethbridge/burn```

* **Method:**

  `POST`

* **Data Params**

  ```
  BurnOrLockReq: { ... }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**

      ```
      { "type", "value" }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Burn given amount of token.




<br />

### Lock tokens

* **URL**

  ```/ethbridge/lock```

* **Method:**

  `POST`

* **Data Params**

  ```
  BurnOrLockReq { ... }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**
      
      ```
      { "type", false }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Lock given amount of token.




<br />

# Clp module




<br />

### Get pool

* **URL**

  ```/clp/getPool```

* **Method:**

  `GET`
  
* **URL Params**

  **Required:**

    `ticker=[string]`

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**
      
      ```
      RawPool{ ... }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Get liquidity pool with a given ticker.




<br />

### Get pools

* **URL**

  ```/clp/getPools```

* **Method:**

  `GET`

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**

      ```
      [ RawPool1, RawPool2, ... ]
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  List all liquidity pools.




<br />

### Create pool

* **URL**

  ```/clp/createPool```

* **Method:**

  `POST`
  
* **Data Params**

  ```
  LiquidityParams: { ... }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**
      
      ```
      { "type", "9000222000444000666" }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Create liquidity pool.




<br />

### Swap tokens

* **URL**

  ```/clp/swap```

* **Method:**

  `POST`
  
* **Data Params**

  ```
  SwapParams: { ... }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**
      
      ```
      { "type", 100 }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Swap sent_asset for received_asset.




<br />

### Get liquidity provider

* **URL**

  ```/clp/getLiquidityProvider```

* **Method:**

  `GET`
  
* **URL Params**

  **Required:**

    ```symbol=[string]```
    ```lpAddress=[string]```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**

      ```
      LiquidityDetailsResponse: { ... }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Get liquidity provider based on address and token symbol.




<br />

### Get liquidity provider assets

* **URL**

  ```/clp/getAssets```

* **Method:**

  `GET`

* **URL Params**

  **Required:**

    ```lpAddress=[string]```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**
      
      ```
      [ "asset1", "asset2", ... ]
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Get assets of specified liquidity provider.




<br />

### Add liquidity

* **URL**

  ```/clp/addLiquidity```

* **Method:**

  `POST`
  
* **Data Params**

  ```
  LiquidityParams: { ... }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**

      ```
      { "type", 123.456 }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Add given amount of token(s) to the liquidity pool.




<br />

### Remove liquidity

* **URL**

  ```/clp/removeLiquidity```

* **Method:**

  `POST`
  
* **Data Params**

  ```
  RemoveLiquidityParams: { .. }
  ```

* **Success Response:**

  * **Code:** 200 <br />

    **Content:**
      
      ```
      { ... }
      ```

* **Error Response:**

  * **Code:** 400 Invalid Request<br />

* **Notes:**

  Remove given amount of token(s) from the liquidity pool.
