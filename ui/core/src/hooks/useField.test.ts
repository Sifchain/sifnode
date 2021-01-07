import { ref, Ref } from "@vue/reactivity";
import { Asset, IAssetAmount, Network, Token } from "../entities";
import { useField } from "./useField";

// TODO eventually delete me as this is an implementation detail
describe("useField", () => {
  let amount: Ref<string>;
  let symbol: Ref<string | null>;

  let asset: Ref<Asset | null>;
  let fieldAmount: Ref<IAssetAmount | null>;
  Token({
    address: "1234",
    decimals: 6,
    name: "Apple Token",
    network: Network.SIFCHAIN,
    symbol: "atk",
  });
  beforeEach(() => {
    amount = ref("0");
    symbol = ref<string | null>(null);
    ({ asset, fieldAmount } = useField(amount, symbol));
  });

  it("should reflect correct values", () => {
    symbol.value = "atk";
    amount.value = "12";

    expect(asset.value?.symbol).toBe("atk");
    expect(fieldAmount.value?.toFixed()).toBe("12.000000");
    amount.value = "123.123123";
    expect(fieldAmount.value?.toFixed()).toBe("123.123123");
  });
});
