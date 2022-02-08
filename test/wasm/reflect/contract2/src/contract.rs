use cosmwasm_std::{entry_point, CosmosMsg};
use cosmwasm_std::{Binary, DepsMut, Env, MessageInfo, Response};

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
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<CustomMsg>, ReflectError> {
    match msg {
        ExecuteMsg::ReflectMsg { msgs } => try_reflect(deps, env, info, msgs),
    }
}

pub fn try_reflect(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msgs: Vec<CosmosMsg<CustomMsg>>,
) -> Result<Response<CustomMsg>, ReflectError> {

    if msgs.is_empty() {
        return Err(ReflectError::MessagesEmpty);
    }

    Ok(Response::new()
        .add_attribute("action", "reflect")
        .add_messages(msgs))
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum CustomMsg {
    Raw(Binary),
}

impl cosmwasm_std::CustomMsg for CustomMsg {}

impl From<CustomMsg> for CosmosMsg<CustomMsg> {
    fn from(original: CustomMsg) -> Self {
        CosmosMsg::Custom(original)
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    ReflectMsg { msgs: Vec<CosmosMsg<CustomMsg>> },
}


#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)] //JsonSchema removed
pub struct InstantiateMsg {}

#[derive(Error, Debug)]
pub enum ReflectError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Messages empty. Must reflect at least one message")]
    MessagesEmpty,
}