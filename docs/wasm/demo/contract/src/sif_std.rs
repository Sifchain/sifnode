use cosmwasm_std::{CosmosMsg, CustomQuery};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum SifchainMsg {
    Swap {
        sent_asset: String,
        received_asset: String,
        sent_amount: String,
        min_received_amount: String,
    },
    AddLiquidity {
        external_asset: String,
        native_asset_amount: String,
        external_asset_amount: String,
    },
    RemoveLiquidity {
        external_asset: String,
        w_basis_points: String,
        asymmetry: String,
    },
}

impl cosmwasm_std::CustomMsg for SifchainMsg {}

impl From<SifchainMsg> for CosmosMsg<SifchainMsg> {
    fn from(original: SifchainMsg) -> Self {
        CosmosMsg::Custom(original)
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum SifchainQuery {
    Pool { external_asset: String },
}

impl CustomQuery for SifchainQuery {}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]

pub struct PoolResponse {
    pub external_asset: String,
    pub external_asset_balance: String,
    pub native_asset_balance: String,
    pub pool_units: String,
}
