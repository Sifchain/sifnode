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
    msg: ReflectCustomMsg,
) -> Result<Response<ReflectCustomMsg>, ReflectError> {    
    Ok(Response::new()
    .add_attribute("action", "reflect")
    .add_message(msg))
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ReflectCustomMsg {
    Swap(String),
}

impl cosmwasm_std::CustomMsg for ReflectCustomMsg {}

impl From<ReflectCustomMsg> for CosmosMsg<ReflectCustomMsg> {
    fn from(original: ReflectCustomMsg) -> Self {
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