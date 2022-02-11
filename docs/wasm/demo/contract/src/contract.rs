use cosmwasm_std::{entry_point, to_binary};
use cosmwasm_std::{Deps, DepsMut, Env, MessageInfo};
use cosmwasm_std::{QueryResponse, Response, StdError, StdResult};
use cosmwasm_std::{Uint256};

use schemars::JsonSchema;
use thiserror::Error;

use serde::{Deserialize, Serialize};

use crate::sif_std::{PoolResponse, SifchainMsg, SifchainQuery};

#[derive(Error, Debug)]
pub enum SwapperError {
    #[error("{0}")]
    Std(#[from] StdError),
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Instantiate
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)] //JsonSchema removed
pub struct InstantiateMsg {}

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, SwapperError> {
    Ok(Response::default())
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Execute
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    Swap {
        amount: u32,
    },
    AddLiquidity {},
    RemoveLiquidity {
        w_basis_points: String,
        asymmetry: String,
    },
}

#[entry_point]
pub fn execute(
    deps: DepsMut<SifchainQuery>,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<SifchainMsg>, SwapperError> {
    match msg {
        ExecuteMsg::Swap { amount } => {

            let pool_response = query_pool(
                deps.as_ref(),
                 "ceth".to_string(),
            )?;

            let external_balance:Uint256 = pool_response.external_asset_balance.parse().unwrap();

            if external_balance < Uint256::from(2_000_000_000_000_000_000u128){
                return Err(SwapperError::Std(StdError::ParseErr{
                    target_type: "xxx".to_string(), 
                    msg: "pool is below threshold".to_string(),
                }))
            }

            let swap_msg = SifchainMsg::Swap {
                sent_asset: "rowan".to_string(),
                received_asset: "ceth".to_string(),
                sent_amount: amount.to_string(),
                min_received_amount: "0".to_string(),
            };

            Ok(Response::new()
                .add_attribute("action", "swap")
                .add_message(swap_msg))
        }
        ExecuteMsg::AddLiquidity {} => {
            let add_liquidity_msg = SifchainMsg::AddLiquidity {
                external_asset: "ceth".to_string(),
                native_asset_amount: "200".to_string(),
                external_asset_amount: "200".to_string(),
            };

            Ok(Response::new()
                .add_attribute("action", "add_liquidity")
                .add_message(add_liquidity_msg))
        }
        ExecuteMsg::RemoveLiquidity {
            w_basis_points,
            asymmetry,
        } => {
            let remove_liquidity_msg = SifchainMsg::RemoveLiquidity {
                external_asset: "ceth".to_string(),
                w_basis_points: w_basis_points,
                asymmetry: asymmetry,
            };

            Ok(Response::new()
                .add_attribute("action", "remove_liquidity")
                .add_message(remove_liquidity_msg))
        }
    }
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Query
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    Pool { external_asset: String },
}

#[entry_point]
pub fn query(deps: Deps<SifchainQuery>, _env: Env, msg: QueryMsg) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::Pool { external_asset } => to_binary(&query_pool(deps, external_asset)?),
    }
}

fn query_pool(deps: Deps<SifchainQuery>, external_asset: String) -> StdResult<PoolResponse> {
    let req = SifchainQuery::Pool { external_asset }.into();
    deps.querier.query(&req)
}
