const { ADMIN_ADDRESS, SIFNODE_CHAIN_ID, SIFNODE_NODE } = process.env;

export const binary = "sifnoded";
export const keyringFlags = ["--keyring-backend=test"];
export const queryFlags = [
  `--node=${SIFNODE_NODE}`,
  `--chain-id=${SIFNODE_CHAIN_ID}`,
];
export const flags = [
  ...keyringFlags,
  ...queryFlags,
  `--from=${ADMIN_ADDRESS}`,
  `-y`,
];
export const userFlags = [...keyringFlags, ...queryFlags, `-y`];
export const gasFlags = [`--gas=500000`, `--gas-prices=0.5rowan`];
export const feesFlags = [`--fees=100000000000000000rowan`];
