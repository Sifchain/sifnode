use cosmwasm_std::{entry_point, DepsMut, Env, MessageInfo, Response};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct InstantiateMsg {}

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, String> {
    Ok(Response::default())
}

use schemars::JsonSchema;
use sif_std::{SifchainMsg, SifchainQuery};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    Swap { amount: u32 },
}

#[entry_point]
pub fn execute(
    _deps: DepsMut<SifchainQuery>,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<SifchainMsg>, String> {
    match msg {
        ExecuteMsg::Swap { amount } => {
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
    }
}

use cosmwasm_std::{to_binary, Deps, QueryResponse, StdResult};
use sif_std::PoolResponse;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    Pool { external_asset: String },
}

#[entry_point]
pub fn query(deps: Deps<SifchainQuery>, _env: Env, msg: QueryMsg) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::Pool { external_asset } => {
            let req = SifchainQuery::Pool { external_asset }.into();
            to_binary(&deps.querier.query::<PoolResponse>(&req)?)
        }
    }
}
