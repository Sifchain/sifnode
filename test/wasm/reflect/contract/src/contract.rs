use cosmwasm_std::{entry_point, CosmosMsg};
use cosmwasm_std::{DepsMut, Env, MessageInfo, Response};

use cosmwasm_std::StdError;
use schemars::JsonSchema;
use thiserror::Error;

use serde::{Deserialize, Serialize};

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ReflectError> {
    Ok(Response::default())
}

#[entry_point]
pub fn execute(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<CustomMsg>, ReflectError> {
    match msg {
        ExecuteMsg::Swap { amount } => Ok(Response::new()
            .add_attribute("action", "reflect")
            .add_message(CustomMsg::Swap { amount })),
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    Swap { amount: u32 },
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum CustomMsg {
    Swap { amount: u32 },
}

impl cosmwasm_std::CustomMsg for CustomMsg {}

impl From<CustomMsg> for CosmosMsg<CustomMsg> {
    fn from(original: CustomMsg) -> Self {
        CosmosMsg::Custom(original)
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)] //JsonSchema removed
pub struct InstantiateMsg {}

#[derive(Error, Debug)]
pub enum ReflectError {
    #[error("{0}")]
    Std(#[from] StdError),
}
