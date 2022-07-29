export function getFee(accounts) {
  const multiplier = accounts.reduce(
    (acc, { pools }) => acc + pools.length + 1,
    0
  );
  const fee = {
    amount: [
      {
        denom: "rowan",
        amount: "10000000000000000000", // 0.1 ROWAN
      },
    ],
    gas: String(18000000 * multiplier), // 180k * multiplier
  };
  return fee;
}
