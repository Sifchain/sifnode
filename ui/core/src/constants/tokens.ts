import { ChainId, Coin, Token } from "../entities";

export const ETH = Coin({
  symbol: "ETH",
  decimals: 18,
  name: "Etherium",
  chainId: ChainId.ETHEREUM,
});

export const ROWAN = Coin({
  symbol: "ROWAN",
  decimals: 2,
  name: "Rowan",
  chainId: ChainId.SIFCHAIN,
});

export const NCN = Coin({
  symbol: "nametoken",
  decimals: 0,
  name: "nametoken",
  chainId: ChainId.SIFCHAIN,
});

export const ZRX = Token({
  symbol: "ZRX",
  decimals: 6,
  name: "0x",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AAVE = Token({
  symbol: "AAVE",
  decimals: 6,
  name: "Aave [New]",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ELF = Token({
  symbol: "ELF",
  decimals: 6,
  name: "elf",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AION = Token({
  symbol: "AION",
  decimals: 6,
  name: "Aion",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AMPL = Token({
  symbol: "AMPL",
  decimals: 6,
  name: "Ampleforth",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ANKR = Token({
  symbol: "ANKR",
  decimals: 6,
  name: "Ankr",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ANT = Token({
  symbol: "ANT",
  decimals: 6,
  name: "Aragon",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAL = Token({
  symbol: "BAL",
  decimals: 6,
  name: "Balancer",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNANA = Token({
  symbol: "BNANA",
  decimals: 6,
  name: "Chimpion",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNT = Token({
  symbol: "BNT",
  decimals: 6,
  name: "Bancor Network Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAND = Token({
  symbol: "BAND",
  decimals: 6,
  name: "Band Protocol",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAT = Token({
  symbol: "BAT",
  decimals: 6,
  name: "Basic Attention Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNB = Token({
  symbol: "BNB",
  decimals: 6,
  name: "Binance Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BUSD = Token({
  symbol: "BUSD",
  decimals: 6,
  name: "Binance USD",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BTMX = Token({
  symbol: "BTMX",
  decimals: 6,
  name: "Bitmax Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BTM = Token({
  symbol: "BTM",
  decimals: 6,
  name: "Bytom",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CEL = Token({
  symbol: "CEL",
  decimals: 6,
  name: "Celsius Network",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LINK = Token({
  symbol: "LINK",
  decimals: 6,
  name: "Chainlink",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CHZ = Token({
  symbol: "CHZ",
  decimals: 6,
  name: "Chiliz",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const COMP = Token({
  symbol: "COMP",
  decimals: 6,
  name: "Compound",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CRO = Token({
  symbol: "CRO",
  decimals: 6,
  name: "Crypto.com Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CRV = Token({
  symbol: "CRV",
  decimals: 6,
  name: "Curve DAO Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const DAI = Token({
  symbol: "DAI",
  decimals: 6,
  name: "Dai",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MANA = Token({
  symbol: "MANA",
  decimals: 6,
  name: "Decentraland",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const DX = Token({
  symbol: "DX",
  decimals: 6,
  name: "DxChain Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ENG = Token({
  symbol: "ENG",
  decimals: 6,
  name: "Enigma",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ENJ = Token({
  symbol: "ENJ",
  decimals: 6,
  name: "Enjin Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LEND = Token({
  symbol: "LEND",
  decimals: 6,
  name: "Aave",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const FTM = Token({
  symbol: "FTM",
  decimals: 6,
  name: "Fantom",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const FET = Token({
  symbol: "FET",
  decimals: 6,
  name: "FirstEnergy Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const GNO = Token({
  symbol: "GNO",
  decimals: 6,
  name: "Gnosis",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const GNT = Token({
  symbol: "GNT",
  decimals: 6,
  name: "Golem",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ONE = Token({
  symbol: "ONE",
  decimals: 6,
  name: "One Hundred Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SNX = Token({
  symbol: "SNX",
  decimals: 6,
  name: "Synthetix Network Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HOT = Token({
  symbol: "HOT",
  decimals: 6,
  name: "Hydro Protocol",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HT = Token({
  symbol: "HT",
  decimals: 6,
  name: "Huobi Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HUSD = Token({
  symbol: "HUSD",
  decimals: 6,
  name: "HUSD",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RLC = Token({
  symbol: "RLC",
  decimals: 6,
  name: "iExec RLC",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const IOST = Token({
  symbol: "IOST",
  decimals: 6,
  name: "IOST",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const IOTX = Token({
  symbol: "IOTX",
  decimals: 6,
  name: "IoTeX",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KEEP = Token({
  symbol: "KEEP",
  decimals: 6,
  name: "Keep Network",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KCS = Token({
  symbol: "KCS",
  decimals: 6,
  name: "KuCoin Shares",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KNC = Token({
  symbol: "KNC",
  decimals: 6,
  name: "Kyber Network",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LEO = Token({
  symbol: "LEO",
  decimals: 6,
  name: "LEO Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LPT = Token({
  symbol: "LPT",
  decimals: 6,
  name: "Livepeer",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LRC = Token({
  symbol: "LRC",
  decimals: 6,
  name: "Loopring",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MKR = Token({
  symbol: "MKR",
  decimals: 6,
  name: "Maker",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MATIC = Token({
  symbol: "MATIC",
  decimals: 6,
  name: "Matic Network",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MCO = Token({
  symbol: "MCO",
  decimals: 6,
  name: "MCO",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MXC = Token({
  symbol: "MXC",
  decimals: 6,
  name: "MXC",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NEXO = Token({
  symbol: "NEXO",
  decimals: 6,
  name: "NEXO",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NMR = Token({
  symbol: "NMR",
  decimals: 6,
  name: "Numeraire",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NXM = Token({
  symbol: "NXM",
  decimals: 6,
  name: "Nexus Mutual",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OCEAN = Token({
  symbol: "OCEAN",
  decimals: 6,
  name: "Ocean Protocol",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OKB = Token({
  symbol: "OKB",
  decimals: 6,
  name: "OKB",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OMG = Token({
  symbol: "OMG",
  decimals: 6,
  name: "OMG Network",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const TRAC = Token({
  symbol: "TRAC",
  decimals: 6,
  name: "OriginTrail",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const PAXG = Token({
  symbol: "PAXG",
  decimals: 6,
  name: "PAX Gold",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const PAX = Token({
  symbol: "PAX",
  decimals: 6,
  name: "PayperEx",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NPXS = Token({
  symbol: "NPXS",
  decimals: 6,
  name: "Pundi X",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const QNT = Token({
  symbol: "QNT",
  decimals: 6,
  name: "Quant",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const REN = Token({
  symbol: "REN",
  decimals: 6,
  name: "REN",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RSR = Token({
  symbol: "RSR",
  decimals: 6,
  name: "Reserve Rights Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RPL = Token({
  symbol: "RPL",
  decimals: 6,
  name: "Rocket Pool",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SRM = Token({
  symbol: "SRM",
  decimals: 6,
  name: "Serum",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AGI = Token({
  symbol: "AGI",
  decimals: 6,
  name: "SingularityNET",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const EURS = Token({
  symbol: "EURS",
  decimals: 6,
  name: "STASIS EURO",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SNT = Token({
  symbol: "SNT",
  decimals: 6,
  name: "Status",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const STORJ = Token({
  symbol: "STORJ",
  decimals: 6,
  name: "Storj",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SUSHI = Token({
  symbol: "SUSHI",
  decimals: 6,
  name: "Sushi",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SXP = Token({
  symbol: "SXP",
  decimals: 6,
  name: "Swipe",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CHSB = Token({
  symbol: "CHSB",
  decimals: 6,
  name: "SwissBorg",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const USDT = Token({
  symbol: "USDT",
  decimals: 6,
  name: "Tether",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const THETA = Token({
  symbol: "THETA",
  decimals: 6,
  name: "Theta Network",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const TUSD = Token({
  symbol: "TUSD",
  decimals: 6,
  name: "TrueUSD",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UMA = Token({
  symbol: "UMA",
  decimals: 6,
  name: "UMA",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UBT = Token({
  symbol: "UBT",
  decimals: 6,
  name: "Unibright",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UNI = Token({
  symbol: "UNI",
  decimals: 6,
  name: "UNIVERSE Token",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UQC = Token({
  symbol: "UQC",
  decimals: 6,
  name: "Uquid Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const USDC = Token({
  symbol: "USDC",
  decimals: 6,
  name: "USD Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UTK = Token({
  symbol: "UTK",
  decimals: 6,
  name: "UTRUST",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const WIC = Token({
  symbol: "WIC",
  decimals: 6,
  name: "Wi Coin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const WBTC = Token({
  symbol: "WBTC",
  decimals: 6,
  name: "Wrapped Bitcoin",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const STAKE = Token({
  symbol: "STAKE",
  decimals: 6,
  name: "xDAI Stake",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const YFI = Token({
  symbol: "YFI",
  decimals: 6,
  name: "yearn.finance",
  chainId: ChainId.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ZIL = Token({
  symbol: "ZIL",
  decimals: 6,
  name: "Zilliqa",
  chainId: ChainId.ETHEREUM,
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
