import { pickChains } from "../pickChains.mjs";

test("pick chains", () => {
  const result = pickChains({ chain: "sifnode,cosmos,akash" });

  expect(result).toMatchSnapshot();
});
