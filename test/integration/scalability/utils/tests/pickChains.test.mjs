import assert from "assert";
import { pickChains } from "../pickChains.mjs";

function test() {
  const result = pickChains({ chain: "sifnode,cosmos,akash" });

  assert.deepEqual(result, ["sifnode", "cosmos", "akash"]);
}

test();
