import { Network, Coin, Token } from "../entities";

export const ETH = Coin({
  symbol: "eth",
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
  symbol: "zrx",
  decimals: 6,
  name: "0x",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AAVE = Token({
  symbol: "aave",
  decimals: 6,
  name: "Aave [New]",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ELF = Token({
  symbol: "elf",
  decimals: 6,
  name: "elf",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AION = Token({
  symbol: "aion",
  decimals: 6,
  name: "Aion",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AMPL = Token({
  symbol: "ampl",
  decimals: 6,
  name: "Ampleforth",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ANKR = Token({
  symbol: "ankr",
  decimals: 6,
  name: "Ankr",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ANT = Token({
  symbol: "ant",
  decimals: 6,
  name: "Aragon",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAL = Token({
  symbol: "bal",
  decimals: 6,
  name: "Balancer",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNANA = Token({
  symbol: "bnana",
  decimals: 6,
  name: "Chimpion",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNT = Token({
  symbol: "bnt",
  decimals: 6,
  name: "Bancor Network Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAND = Token({
  symbol: "band",
  decimals: 6,
  name: "Band Protocol",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BAT = Token({
  symbol: "bat",
  decimals: 6,
  name: "Basic Attention Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BNB = Token({
  symbol: "bnb",
  decimals: 6,
  name: "Binance Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BUSD = Token({
  symbol: "busd",
  decimals: 6,
  name: "Binance USD",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BTMX = Token({
  symbol: "btmx",
  decimals: 6,
  name: "Bitmax Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const BTM = Token({
  symbol: "btm",
  decimals: 6,
  name: "Bytom",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CEL = Token({
  symbol: "cel",
  decimals: 6,
  name: "Celsius Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LINK = Token({
  symbol: "link",
  decimals: 6,
  name: "Chainlink",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CHZ = Token({
  symbol: "chz",
  decimals: 6,
  name: "Chiliz",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const COMP = Token({
  symbol: "comp",
  decimals: 6,
  name: "Compound",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CRO = Token({
  symbol: "cro",
  decimals: 6,
  name: "Crypto.com Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CRV = Token({
  symbol: "crv",
  decimals: 6,
  name: "Curve DAO Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const DAI = Token({
  symbol: "dai",
  decimals: 6,
  name: "Dai",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MANA = Token({
  symbol: "mana",
  decimals: 6,
  name: "Decentraland",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const DX = Token({
  symbol: "dx",
  decimals: 6,
  name: "DxChain Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ENG = Token({
  symbol: "eng",
  decimals: 6,
  name: "Enigma",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ENJ = Token({
  symbol: "enj",
  decimals: 6,
  name: "Enjin Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LEND = Token({
  symbol: "lend",
  decimals: 6,
  name: "Aave",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const FTM = Token({
  symbol: "ftm",
  decimals: 6,
  name: "Fantom",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const FET = Token({
  symbol: "fet",
  decimals: 6,
  name: "FirstEnergy Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const GNO = Token({
  symbol: "gno",
  decimals: 6,
  name: "Gnosis",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const GNT = Token({
  symbol: "gnt",
  decimals: 6,
  name: "Golem",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ONE = Token({
  symbol: "one",
  decimals: 6,
  name: "One Hundred Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SNX = Token({
  symbol: "snx",
  decimals: 6,
  name: "Synthetix Network Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HOT = Token({
  symbol: "hot",
  decimals: 6,
  name: "Hydro Protocol",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HT = Token({
  symbol: "ht",
  decimals: 6,
  name: "Huobi Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const HUSD = Token({
  symbol: "husd",
  decimals: 6,
  name: "HUSD",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RLC = Token({
  symbol: "rlc",
  decimals: 6,
  name: "iExec RLC",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const IOST = Token({
  symbol: "iost",
  decimals: 6,
  name: "IOST",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const IOTX = Token({
  symbol: "iotx",
  decimals: 6,
  name: "IoTeX",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KEEP = Token({
  symbol: "keep",
  decimals: 6,
  name: "Keep Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KCS = Token({
  symbol: "kcs",
  decimals: 6,
  name: "KuCoin Shares",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const KNC = Token({
  symbol: "knc",
  decimals: 6,
  name: "Kyber Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LEO = Token({
  symbol: "leo",
  decimals: 6,
  name: "LEO Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LPT = Token({
  symbol: "lpt",
  decimals: 6,
  name: "Livepeer",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const LRC = Token({
  symbol: "lrc",
  decimals: 6,
  name: "Loopring",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MKR = Token({
  symbol: "mkr",
  decimals: 6,
  name: "Maker",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MATIC = Token({
  symbol: "matic",
  decimals: 6,
  name: "Matic Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MCO = Token({
  symbol: "mco",
  decimals: 6,
  name: "MCO",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const MXC = Token({
  symbol: "mxc",
  decimals: 6,
  name: "MXC",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NEXO = Token({
  symbol: "nexo",
  decimals: 6,
  name: "NEXO",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NMR = Token({
  symbol: "nmr",
  decimals: 6,
  name: "Numeraire",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NXM = Token({
  symbol: "nxm",
  decimals: 6,
  name: "Nexus Mutual",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OCEAN = Token({
  symbol: "ocean",
  decimals: 6,
  name: "Ocean Protocol",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OKB = Token({
  symbol: "okb",
  decimals: 6,
  name: "OKB",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const OMG = Token({
  symbol: "omg",
  decimals: 6,
  name: "OMG Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const TRAC = Token({
  symbol: "trac",
  decimals: 6,
  name: "OriginTrail",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const PAXG = Token({
  symbol: "paxg",
  decimals: 6,
  name: "PAX Gold",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const PAX = Token({
  symbol: "pax",
  decimals: 6,
  name: "PayperEx",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const NPXS = Token({
  symbol: "npxs",
  decimals: 6,
  name: "Pundi X",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const QNT = Token({
  symbol: "qnt",
  decimals: 6,
  name: "Quant",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const REN = Token({
  symbol: "ren",
  decimals: 6,
  name: "REN",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RSR = Token({
  symbol: "rsr",
  decimals: 6,
  name: "Reserve Rights Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const RPL = Token({
  symbol: "rpl",
  decimals: 6,
  name: "Rocket Pool",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SRM = Token({
  symbol: "srm",
  decimals: 6,
  name: "Serum",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const AGI = Token({
  symbol: "agi",
  decimals: 6,
  name: "SingularityNET",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const EURS = Token({
  symbol: "eurs",
  decimals: 6,
  name: "STASIS EURO",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SNT = Token({
  symbol: "snt",
  decimals: 6,
  name: "Status",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const STORJ = Token({
  symbol: "storj",
  decimals: 6,
  name: "Storj",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SUSHI = Token({
  symbol: "sushi",
  decimals: 6,
  name: "Sushi",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const SXP = Token({
  symbol: "sxp",
  decimals: 6,
  name: "Swipe",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const CHSB = Token({
  symbol: "chsb",
  decimals: 6,
  name: "SwissBorg",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const USDT = Token({
  symbol: "usdt",
  decimals: 6,
  name: "Tether",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const THETA = Token({
  symbol: "theta",
  decimals: 6,
  name: "Theta Network",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const TUSD = Token({
  symbol: "tusd",
  decimals: 6,
  name: "TrueUSD",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UMA = Token({
  symbol: "uma",
  decimals: 6,
  name: "UMA",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UBT = Token({
  symbol: "ubt",
  decimals: 6,
  name: "Unibright",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UNI = Token({
  symbol: "uni",
  decimals: 6,
  name: "UNIVERSE Token",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UQC = Token({
  symbol: "uqc",
  decimals: 6,
  name: "Uquid Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const USDC = Token({
  symbol: "usdc",
  decimals: 6,
  name: "USD Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const UTK = Token({
  symbol: "utk",
  decimals: 6,
  name: "UTRUST",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const WIC = Token({
  symbol: "wic",
  decimals: 6,
  name: "Wi Coin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const WBTC = Token({
  symbol: "wbtc",
  decimals: 6,
  name: "Wrapped Bitcoin",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const STAKE = Token({
  symbol: "stake",
  decimals: 6,
  name: "xDAI Stake",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const YFI = Token({
  symbol: "yfi",
  decimals: 6,
  name: "yearn.finance",
  network: Network.ETHEREUM,
  address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", //USDC Contract address TODO - replacewith real one
});

export const ZIL = Token({
  symbol: "zil",
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
