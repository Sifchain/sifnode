import { ChainId, createCoin, createToken } from "../entities";

// This should all probably be relocated to a token service a JSON or built off proper data sources

export const ETH = createCoin("ETH", 18, "Etherium", ChainId.ETHEREUM);
export const ROWAN = createCoin("ROWAN", 2, "Rowan", ChainId.SIFCHAIN);
export const NCN = createCoin("nametoken", 0, "nametoken", ChainId.SIFCHAIN);

export const ZRX = createToken(
  "ZRX",
  6,
  "0x",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const AAVE = createToken(
  "AAVE",
  6,
  "Aave [New]",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ELF = createToken(
  "ELF",
  6,
  "elf",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const AION = createToken(
  "AION",
  6,
  "Aion",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const AMPL = createToken(
  "AMPL",
  6,
  "Ampleforth",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ANKR = createToken(
  "ANKR",
  6,
  "Ankr",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ANT = createToken(
  "ANT",
  6,
  "Aragon",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BAL = createToken(
  "BAL",
  6,
  "Balancer",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BNANA = createToken(
  "BNANA",
  6,
  "Chimpion",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BNT = createToken(
  "BNT",
  6,
  "Bancor Network Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BAND = createToken(
  "BAND",
  6,
  "Band Protocol",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BAT = createToken(
  "BAT",
  6,
  "Basic Attention Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BNB = createToken(
  "BNB",
  6,
  "Binance Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BUSD = createToken(
  "BUSD",
  6,
  "Binance USD",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BTMX = createToken(
  "BTMX",
  6,
  "Bitmax Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const BTM = createToken(
  "BTM",
  6,
  "Bytom",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const CEL = createToken(
  "CEL",
  6,
  "Celsius Network",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const LINK = createToken(
  "LINK",
  6,
  "Chainlink",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const CHZ = createToken(
  "CHZ",
  6,
  "Chiliz",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const COMP = createToken(
  "COMP",
  6,
  "Compound",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const CRO = createToken(
  "CRO",
  6,
  "Crypto.com Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const CRV = createToken(
  "CRV",
  6,
  "Curve DAO Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const DAI = createToken(
  "DAI",
  6,
  "Dai",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const MANA = createToken(
  "MANA",
  6,
  "Decentraland",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const DX = createToken(
  "DX",
  6,
  "DxChain Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ENG = createToken(
  "ENG",
  6,
  "Enigma",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ENJ = createToken(
  "ENJ",
  6,
  "Enjin Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const LEND = createToken(
  "LEND",
  6,
  "Aave",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const FTM = createToken(
  "FTM",
  6,
  "Fantom",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const FET = createToken(
  "FET",
  6,
  "FirstEnergy Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const GNO = createToken(
  "GNO",
  6,
  "Gnosis",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const GNT = createToken(
  "GNT",
  6,
  "Golem",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ONE = createToken(
  "ONE",
  6,
  "One Hundred Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const SNX = createToken(
  "SNX",
  6,
  "Synthetix Network Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const HOT = createToken(
  "HOT",
  6,
  "Hydro Protocol",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const HT = createToken(
  "HT",
  6,
  "Huobi Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const HUSD = createToken(
  "HUSD",
  6,
  "HUSD",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const RLC = createToken(
  "RLC",
  6,
  "iExec RLC",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const IOST = createToken(
  "IOST",
  6,
  "IOST",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const IOTX = createToken(
  "IOTX",
  6,
  "IoTeX",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const KEEP = createToken(
  "KEEP",
  6,
  "Keep Network",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const KCS = createToken(
  "KCS",
  6,
  "KuCoin Shares",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const KNC = createToken(
  "KNC",
  6,
  "Kyber Network",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const LEO = createToken(
  "LEO",
  6,
  "LEO Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const LPT = createToken(
  "LPT",
  6,
  "Livepeer",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const LRC = createToken(
  "LRC",
  6,
  "Loopring",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const MKR = createToken(
  "MKR",
  6,
  "Maker",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const MATIC = createToken(
  "MATIC",
  6,
  "Matic Network",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const MCO = createToken(
  "MCO",
  6,
  "MCO",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const MXC = createToken(
  "MXC",
  6,
  "MXC",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const NEXO = createToken(
  "NEXO",
  6,
  "NEXO",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const NMR = createToken(
  "NMR",
  6,
  "Numeraire",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const NXM = createToken(
  "NXM",
  6,
  "Nexus Mutual",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const OCEAN = createToken(
  "OCEAN",
  6,
  "Ocean Protocol",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const OKB = createToken(
  "OKB",
  6,
  "OKB",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const OMG = createToken(
  "OMG",
  6,
  "OMG Network",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const TRAC = createToken(
  "TRAC",
  6,
  "OriginTrail",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const PAXG = createToken(
  "PAXG",
  6,
  "PAX Gold",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const PAX = createToken(
  "PAX",
  6,
  "PayperEx",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const NPXS = createToken(
  "NPXS",
  6,
  "Pundi X",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const QNT = createToken(
  "QNT",
  6,
  "Quant",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const REN = createToken(
  "REN",
  6,
  "REN",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const RSR = createToken(
  "RSR",
  6,
  "Reserve Rights Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const RPL = createToken(
  "RPL",
  6,
  "Rocket Pool",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const SRM = createToken(
  "SRM",
  6,
  "Serum",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const AGI = createToken(
  "AGI",
  6,
  "SingularityNET",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const EURS = createToken(
  "EURS",
  6,
  "STASIS EURO",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const SNT = createToken(
  "SNT",
  6,
  "Status",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const STORJ = createToken(
  "STORJ",
  6,
  "Storj",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const SUSHI = createToken(
  "SUSHI",
  6,
  "Sushi",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const SXP = createToken(
  "SXP",
  6,
  "Swipe",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const CHSB = createToken(
  "CHSB",
  6,
  "SwissBorg",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const USDT = createToken(
  "USDT",
  6,
  "Tether",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const THETA = createToken(
  "THETA",
  6,
  "Theta Network",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const TUSD = createToken(
  "TUSD",
  6,
  "TrueUSD",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const UMA = createToken(
  "UMA",
  6,
  "UMA",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const UBT = createToken(
  "UBT",
  6,
  "Unibright",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const UNI = createToken(
  "UNI",
  6,
  "UNIVERSE Token",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const UQC = createToken(
  "UQC",
  6,
  "Uquid Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const USDC = createToken(
  "USDC",
  6,
  "USD Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const UTK = createToken(
  "UTK",
  6,
  "UTRUST",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const WIC = createToken(
  "WIC",
  6,
  "Wi Coin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const WBTC = createToken(
  "WBTC",
  6,
  "Wrapped Bitcoin",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const STAKE = createToken(
  "STAKE",
  6,
  "xDAI Stake",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const YFI = createToken(
  "YFI",
  6,
  "yearn.finance",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

export const ZIL = createToken(
  "ZIL",
  6,
  "Zilliqa",
  ChainId.ETHEREUM,
  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" //USDC Contract address TODO - replacewith real one
);

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
