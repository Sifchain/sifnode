// see https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md

// calculates liquidity fee per Thorchain's CLP model
export function calcLiquidityFee(
  X: number, // Sent Asset Pool Balance,
  x: number, // Sent Asset Amount,
  Y: number // Received Asset Pool Balance
) {
  return (x * x * Y) / ((x + X) * (x + X));
}

// calculates trade slip per Thorchain's CLP model
export function calcTradeSlip(
  X: number, // Sent Asset Pool Balance,
  x: number // Sent Asset Amount,
) {
  return (x * (2 * X + x)) / (X * X);
}

// calculates final swap received token amount
export function calcSwapResult(
  X: number, // Sent Asset Pool Balance,
  x: number, // Sent Asset Amount,
  Y: number // Received Asset Pool Balance
) {
  return (x * X * Y) / ((x + X) * (x + X));
}
