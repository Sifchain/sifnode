export default ({ api, store }) => ({
    async connectToEthWallet(ethWallet) {
        //
        // This was from the Swap tab
        // GET: <ETH_WALLET>	LocalWallet: <ETH_WALLET>: UserWallet: Info: Object: [Address: String, Connected: Bool] -> LocalStorage  -> LocalStorage(App) (probably web3.isConnected)
        // GET: Tokens In <ETH_WALLET> 	LocalWallet: <ETH_WALLET>: UserWallet: Tokens: Object: [Quantity: Number, Value: Number, Ticker: String, Address: String] -> LocalStorage(App) (probably web3.isConnected)
        //
        // This was from the wallet tab in the spreadsheet
        // GET: <ETH_WALLET>	LocalWallet: <ETH_WALLET>: Object: [Address: String, Connected: Bool] -> LocalStorage
        // GET: Tokens In Wallet 	LocalWallet: UserTokens: Object: [Quantity: Number, Value: Number, Ticker: String, Address: String] -> LocalStorage
    },
    async connectToCosmosWallet(cosmosWallet) {
        // GET: <C_WALLET>	LocalWallet: <C_WALLET>: Object: [Address: String, Connected: Bool] -> LocalStorage
        // GET: Tokens In Wallet 	LocalWallet: UserTokens: Object: [Quantity: Number, Value: Number, Ticker: String, Address: String] -> LocalStorage
    },
});
//# sourceMappingURL=connectToWallet.js.map