// import AssetItemLarge from "./AssetItemLarge.vue";
// import { parseAssets } from "../../../../../core/src/utils/parseConfig";
// import localethereumassets from "../../../../../core/src/assets.ethereum.localnet.json";
// import localsifassets from "../../../../../core/src/assets.sifchain.localnet.json";
// import { Asset, IAssetAmount } from "../../../../../core/src/entities";
// const assets = [...localethereumassets.assets, ...localsifassets.assets];
//
// //core / src / test / utils / getTestingToken.ts;
// // {
// //   "name": "Ethereum",
// //   "symbol": "eth",
// //   "imageUrl": "https://assets.coingecko.com/coins/images/279/large/ethereum.png",
// //   "network": "ethereum",
// //   "decimals": 18
// // },
// export default {
//   title: "AssetItemLarge",
//   component: AssetItemLarge,
// };
//
// const Template = (args: any) => ({
//   props: [],
//   components: { AssetItemLarge },
//   setup() {
//     // const supportedTokens = parseAssets(assets as any[]).map((asset) => {
//     //   Asset.set(asset.symbol, asset);
//     //   return asset;
//     // });
//     const ethAsset = Asset.set("eth", {
//       name: "Ethereum",
//       symbol: "eth",
//       network: "ethereum",
//       decimals: 18,
//     });
//     return { args };
//   },
//   template: '<div><AssetItemLarge v-bind="args" /></div>',
// });
//
// export const Primary = Template.bind({});
// (Primary as any).args = {
//   amount: 1000000,
//   symbol: "eth",
// };
