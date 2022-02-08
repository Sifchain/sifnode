use cosmwasm_std::{entry_point, CosmosMsg};
use cosmwasm_std::{DepsMut, Env, MessageInfo, Response};

use cosmwasm_std::StdError;
use schemars::JsonSchema;
use thiserror::Error;

use serde::{Deserialize, Serialize};

const RAW_MSG: &str = "eyJAdHlwZSI6Ii9zaWZub2RlLmNscC52MS5Nc2dTd2FwIiwic2lnbmVyIjoic2lmMTRoajJ0YXZxOGZwZXNkd3h4Y3U0NHJ0eTNoaDkwdmh1anJ2Y21zdGw0enIzdHhtZnZ3OXM2MmN2dTYiLCJzZW50X2Fzc2V0Ijp7InN5bWJvbCI6InJvd2FuIn0sInJlY2VpdmVkX2Fzc2V0Ijp7InN5bWJvbCI6ImNldGgifSwic2VudF9hbW91bnQiOiIyMDAwMCIsIm1pbl9yZWNlaXZpbmdfYW1vdW50IjoiMCJ9";

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
        ExecuteMsg::Swap {} => {
        
            Ok(Response::new()
                .add_attribute("action", "reflect")
                .add_message( CustomMsg::Raw(RAW_MSG.to_string())))
        }
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum CustomMsg {
    Raw(String),
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
    Swap {},
}


#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)] //JsonSchema removed
pub struct InstantiateMsg {}

#[derive(Error, Debug)]
pub enum ReflectError {
    #[error("{0}")]
    Std(#[from] StdError),
}