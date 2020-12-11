import { Network, Coin, Token } from "../entities";

export const ETH = Coin({
  symbol: "ETH",
  decimals: 18,
  name: "Etherium",
  network: Network.ETHEREUM,
});

export const ROWAN = Coin({
  symbol: "rowan",
  decimals: 18,
  name: "Rowan",
  network: Network.SIFCHAIN,
});

export const ATK = Token({
  symbol: "atk",
  address: "0xbaAA2a3237035A2c7fA2A33c76B44a8C6Fe18e87",
  decimals: 18,
  name: "atk",
  network: Network.ETHEREUM,
});

export const CATK = Coin({
  symbol: "catk",
  decimals: 18,
  name: "catk",
  network: Network.SIFCHAIN,
});

export const CBTK = Coin({
  symbol: "cbtk",
  decimals: 18,
  name: "cbtk",
  network: Network.SIFCHAIN,
});

export const CETH = Coin({
  symbol: "ceth",
  decimals: 18,
  name: "ceth",
  network: Network.SIFCHAIN,
});

export const ZRX = Token({
  symbol: "ZRX",
  decimals: 6,
  name: "0x",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AAVE = Token({
  symbol: "AAVE",
  decimals: 6,
  name: "Aave [New]",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ELF = Token({
  symbol: "ELF",
  decimals: 6,
  name: "elf",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AION = Token({
  symbol: "AION",
  decimals: 6,
  name: "Aion",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AMPL = Token({
  symbol: "AMPL",
  decimals: 6,
  name: "Ampleforth",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ANKR = Token({
  symbol: "ANKR",
  decimals: 6,
  name: "Ankr",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ANT = Token({
  symbol: "ANT",
  decimals: 6,
  name: "Aragon",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAL = Token({
  symbol: "BAL",
  decimals: 6,
  name: "Balancer",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNANA = Token({
  symbol: "BNANA",
  decimals: 6,
  name: "Chimpion",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNT = Token({
  symbol: "BNT",
  decimals: 6,
  name: "Bancor Network Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAND = Token({
  symbol: "BAND",
  decimals: 6,
  name: "Band Protocol",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAT = Token({
  symbol: "BAT",
  decimals: 6,
  name: "Basic Attention Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNB = Token({
  symbol: "BNB",
  decimals: 6,
  name: "Binance Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BUSD = Token({
  symbol: "BUSD",
  decimals: 6,
  name: "Binance USD",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BTMX = Token({
  symbol: "BTMX",
  decimals: 6,
  name: "Bitmax Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BTM = Token({
  symbol: "BTM",
  decimals: 6,
  name: "Bytom",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CEL = Token({
  symbol: "CEL",
  decimals: 6,
  name: "Celsius Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LINK = Token({
  symbol: "LINK",
  decimals: 6,
  name: "Chainlink",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CHZ = Token({
  symbol: "CHZ",
  decimals: 6,
  name: "Chiliz",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const COMP = Token({
  symbol: "COMP",
  decimals: 6,
  name: "Compound",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CRO = Token({
  symbol: "CRO",
  decimals: 6,
  name: "Crypto.com Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CRV = Token({
  symbol: "CRV",
  decimals: 6,
  name: "Curve DAO Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const DAI = Token({
  symbol: "DAI",
  decimals: 6,
  name: "Dai",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MANA = Token({
  symbol: "MANA",
  decimals: 6,
  name: "Decentraland",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const DX = Token({
  symbol: "DX",
  decimals: 6,
  name: "DxChain Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ENG = Token({
  symbol: "ENG",
  decimals: 6,
  name: "Enigma",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ENJ = Token({
  symbol: "ENJ",
  decimals: 6,
  name: "Enjin Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LEND = Token({
  symbol: "LEND",
  decimals: 6,
  name: "Aave",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const FTM = Token({
  symbol: "FTM",
  decimals: 6,
  name: "Fantom",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const FET = Token({
  symbol: "FET",
  decimals: 6,
  name: "FirstEnergy Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const GNO = Token({
  symbol: "GNO",
  decimals: 6,
  name: "Gnosis",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const GNT = Token({
  symbol: "GNT",
  decimals: 6,
  name: "Golem",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ONE = Token({
  symbol: "ONE",
  decimals: 6,
  name: "One Hundred Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SNX = Token({
  symbol: "SNX",
  decimals: 6,
  name: "Synthetix Network Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HOT = Token({
  symbol: "HOT",
  decimals: 6,
  name: "Hydro Protocol",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HT = Token({
  symbol: "HT",
  decimals: 6,
  name: "Huobi Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HUSD = Token({
  symbol: "HUSD",
  decimals: 6,
  name: "HUSD",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RLC = Token({
  symbol: "RLC",
  decimals: 6,
  name: "iExec RLC",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const IOST = Token({
  symbol: "IOST",
  decimals: 6,
  name: "IOST",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const IOTX = Token({
  symbol: "IOTX",
  decimals: 6,
  name: "IoTeX",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KEEP = Token({
  symbol: "KEEP",
  decimals: 6,
  name: "Keep Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KCS = Token({
  symbol: "KCS",
  decimals: 6,
  name: "KuCoin Shares",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KNC = Token({
  symbol: "KNC",
  decimals: 6,
  name: "Kyber Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LEO = Token({
  symbol: "LEO",
  decimals: 6,
  name: "LEO Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LPT = Token({
  symbol: "LPT",
  decimals: 6,
  name: "Livepeer",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LRC = Token({
  symbol: "LRC",
  decimals: 6,
  name: "Loopring",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MKR = Token({
  symbol: "MKR",
  decimals: 6,
  name: "Maker",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MATIC = Token({
  symbol: "MATIC",
  decimals: 6,
  name: "Matic Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MCO = Token({
  symbol: "MCO",
  decimals: 6,
  name: "MCO",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MXC = Token({
  symbol: "MXC",
  decimals: 6,
  name: "MXC",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NEXO = Token({
  symbol: "NEXO",
  decimals: 6,
  name: "NEXO",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NMR = Token({
  symbol: "NMR",
  decimals: 6,
  name: "Numeraire",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NXM = Token({
  symbol: "NXM",
  decimals: 6,
  name: "Nexus Mutual",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OCEAN = Token({
  symbol: "OCEAN",
  decimals: 6,
  name: "Ocean Protocol",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OKB = Token({
  symbol: "OKB",
  decimals: 6,
  name: "OKB",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OMG = Token({
  symbol: "OMG",
  decimals: 6,
  name: "OMG Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const TRAC = Token({
  symbol: "TRAC",
  decimals: 6,
  name: "OriginTrail",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const PAXG = Token({
  symbol: "PAXG",
  decimals: 6,
  name: "PAX Gold",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const PAX = Token({
  symbol: "PAX",
  decimals: 6,
  name: "PayperEx",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NPXS = Token({
  symbol: "NPXS",
  decimals: 6,
  name: "Pundi X",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const QNT = Token({
  symbol: "QNT",
  decimals: 6,
  name: "Quant",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const REN = Token({
  symbol: "REN",
  decimals: 6,
  name: "REN",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RSR = Token({
  symbol: "RSR",
  decimals: 6,
  name: "Reserve Rights Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RPL = Token({
  symbol: "RPL",
  decimals: 6,
  name: "Rocket Pool",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SRM = Token({
  symbol: "SRM",
  decimals: 6,
  name: "Serum",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AGI = Token({
  symbol: "AGI",
  decimals: 6,
  name: "SingularityNET",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const EURS = Token({
  symbol: "EURS",
  decimals: 6,
  name: "STASIS EURO",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SNT = Token({
  symbol: "SNT",
  decimals: 6,
  name: "Status",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const STORJ = Token({
  symbol: "STORJ",
  decimals: 6,
  name: "Storj",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SUSHI = Token({
  symbol: "SUSHI",
  decimals: 6,
  name: "Sushi",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SXP = Token({
  symbol: "SXP",
  decimals: 6,
  name: "Swipe",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CHSB = Token({
  symbol: "CHSB",
  decimals: 6,
  name: "SwissBorg",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const USDT = Token({
  symbol: "USDT",
  decimals: 6,
  name: "Tether",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const THETA = Token({
  symbol: "THETA",
  decimals: 6,
  name: "Theta Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const TUSD = Token({
  symbol: "TUSD",
  decimals: 6,
  name: "TrueUSD",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UMA = Token({
  symbol: "UMA",
  decimals: 6,
  name: "UMA",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UBT = Token({
  symbol: "UBT",
  decimals: 6,
  name: "Unibright",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UNI = Token({
  symbol: "UNI",
  decimals: 6,
  name: "UNIVERSE Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UQC = Token({
  symbol: "UQC",
  decimals: 6,
  name: "Uquid Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const USDC = Token({
  symbol: "USDC",
  decimals: 6,
  name: "USD Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UTK = Token({
  symbol: "UTK",
  decimals: 6,
  name: "UTRUST",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const WIC = Token({
  symbol: "WIC",
  decimals: 6,
  name: "Wi Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const WBTC = Token({
  symbol: "WBTC",
  decimals: 6,
  name: "Wrapped Bitcoin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const STAKE = Token({
  symbol: "STAKE",
  decimals: 6,
  name: "xDAI Stake",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const YFI = Token({
  symbol: "YFI",
  decimals: 6,
  name: "yearn.finance",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ZIL = Token({
  symbol: "ZIL",
  decimals: 6,
  name: "Zilliqa",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

// Here for Reference: dummy arbitrary list replace with dynamic list? https://www.finder.com.au/erc20-tokens
// export const MARKETCAP_TOKEN_ORDER = [
//   "USDT",
//   "BNB",
//   "LINK",
//   "CRO",
//   "USDC",
//   "OKB",
//   "cDAI",
//   "LEO",
//   "HT",
//   "WBTC",
//   "DAI",
//   "VEN",
//   "THETA",
//   "BUSD",
//   "UNI",
//   "LEND",
//   "YFI",
//   "MKR",
//   "SNX",
//   "OMG",
//   "CEL",
//   "UMA",
//   "TUSD",
//   "COMP",
//   "BAT",
//   "cETH",
//   "PAX",
//   "renBTC",
//   "ZRX",
//   "REN",
//   "cUSDC",
//   "AAVE",
//   "NXM",
//   "LRC",
//   "ZIL",
//   "KNC",
//   "NMR",
//   "ENJ",
//   "BAND",
//   "HUSD",
//   "BAL",
//   "ANT",
//   "OCEAN",
//   "QNT",
//   "BTM",
//   "MANA",
//   "SXP",
//   "DX",
//   "SNT",
//   "GNT",
//   "AMPL",
//   "IOST",
//   "HOT",
//   "RSR",
//   "DIVX",
//   "SUSHI",
//   "FTM",
//   "BNT",
//   "KCS",
//   "NEXO",
//   "REPv2",
//   "STORJ",
//   "sUSD",
//   "MCO",
//   "PAXG",
//   "wNXM",
//   "SRM",
//   "MATIC",
//   "cUSDT",
//   "KEEP",
//   "LPT",
//   "CHZ",
//   "UTK",
//   "WAX",
//   "CHSB",
//   "RLC",
//   "cZRX",
//   "MXC",
//   "ENG",
//   "UQC",
//   "CRV",
//   "UBT",
//   "NPXS",
//   "GNO",
//   "ELF",
//   "WIC",
//   "IOTX",
//   "AGI",
//   "RPL",
//   "ANKR",
//   "ONE",
//   "mUSD",
//   "EURS",
//   "STAKE",
//   "AION",
//   "FET",
//   "BTMX",
//   "AURA",
//   "TRAC",
//   "BNANA",
// ];
