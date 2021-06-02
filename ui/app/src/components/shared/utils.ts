import { computed, Ref, ComputedRef } from "@vue/reactivity";
import ColorHash from "color-hash";
import { Asset, IAssetAmount, Network, toBaseUnits, TxHash } from "ui-core";
import { format } from "ui-core/src/utils/format";

export function formatSymbol(symbol: string) {
  if (symbol.indexOf("c") === 0) {
    return ["c", symbol.slice(1).toUpperCase()].join("");
  }
  return symbol.toUpperCase();
}

export function formatPercentage(amount: string) {
  return parseFloat(amount) < 0.01
    ? "< 0.01%"
    : `${parseFloat(amount).toFixed(2)}%`;
}
// TODO: make this work for AssetAmounts and Fractions / Amounts
export function formatNumber(displayNumber: string) {
  if (!displayNumber) return "0";
  const amount = parseFloat(displayNumber);
  if (amount < 100000) {
    return amount.toFixed(6);
  } else {
    return amount.toFixed(2);
  }
}

export function formatAssetAmount(value: IAssetAmount) {
  if (!value || value.equalTo("0")) return "0";
  const { amount, asset } = value;
  return amount.greaterThan(toBaseUnits("100000", asset))
    ? format(amount, asset, { mantissa: 2 })
    : format(amount, asset, { mantissa: 6 });
}

// TODO: These could be replaced with a look up table
export function getPeggedSymbol(symbol: string) {
  if (symbol.toLowerCase() === "erowan") return "ROWAN";
  return "c" + symbol.toUpperCase();
}
export function getUnpeggedSymbol(symbol: string) {
  if (symbol.toLowerCase() === "rowan") return "eROWAN";
  return symbol.indexOf("c") === 0 ? symbol.slice(1) : symbol;
}

export function getAssetLabel(t: Asset) {
  if (t.network === Network.SIFCHAIN) {
    return formatSymbol(t.symbol);
  }

  if (t.network === Network.ETHEREUM && t.symbol.toLowerCase() === "erowan") {
    return "eROWAN";
  }

  return t.symbol.toUpperCase();
}

export function useAssetItem(symbol: Ref<string | undefined>) {
  const token = computed(() =>
    symbol.value ? Asset.get(symbol.value) : undefined,
  );

  const tokenLabel = computed(() => {
    if (!token.value) return "";
    return getAssetLabel(token.value);
  });

  const backgroundStyle = computed(() => {
    if (!symbol.value) return "";

    const colorHash = new ColorHash();

    const color = symbol ? colorHash.hex(symbol.value) : [];

    return `background: ${color};`;
  });

  return {
    token: token,
    label: tokenLabel,
    background: backgroundStyle,
  };
}

export async function getLMData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  // const ceUrl = getCryptoeconomicsUrl(chainId);
  // const data = await fetch(
  //   `${ceUrl}/lm/?key=userData&address=${address.value}&timestamp=now`,
  // );
  // if (data.status !== 200) return {};
  const parsedData = {
    totalDepositedAmount: 53919106.11780828,
    timestamp: 147400,
    rewardBuckets: [
      {
        rowan: 17339449.541285392,
        initialRowan: 45000000,
        duration: 1200,
      },
    ],
    user: {
      tickets: [
        {
          commission: 0,
          amount: 293.99999564707366,
          mul: 0.417534722222227,
          reward: 38.12101358818803,
          validatorRewardAddress: null,
          validatorStakeAddress: null,
          timestamp: "May 5th 2021, 8:08:43 pm",
          rewardDelta: 0.20469172665896876,
          poolDominanceRatio: 0.000005453897339202301,
          commissionRewardsByValidator: {},
        },
      ],
      claimableRewardsOnWithdrawnAssets: 0,
      dispensed: 0,
      forfeited: 0,
      totalAccruedCommissionsAndClaimableRewards: 15.916846819373829,
      totalClaimableCommissionsAndClaimableRewards: 15.916846819373829,
      reservedReward: 38.12101358818803,
      totalDepositedAmount: 293.99999564707366,
      totalClaimableRewardsOnDepositedAssets: 15.916846819373829,
      currentTotalCommissionsOnClaimableDelegatorRewards: 0,
      totalAccruedCommissionsAtMaturity: 0,
      totalCommissionsAndRewardsAtMaturity: 132.86652090850077,
      claimableCommissions: 0,
      delegatorAddresses: [],
      totalRewardsOnDepositedAssetsAtMaturity: 132.86652090850077,
      ticketAmountAtMaturity: 293.99999564707366,
      yieldAtMaturity: 0.4519269485568214,
      nextRewardShare: 0.000005452612567513837,
      currentYieldOnTickets: 0.39778801299547234,
      maturityDate: "September 2nd 2021, 8:08:43 pm",
      maturityDateISO: "2021-09-02T20:08:43.000Z",
      yearsToMaturity: 0.2553272450532724,
      currentAPYOnTickets: 1.557953648512819,
      maturityDateMs: 0,
      futureReward: 116.94967408912694,
      nextReward: 0.2044729712817689,
      nextRewardProjectedFutureReward: 537.3549685284886,
      nextRewardProjectedAPYOnTickets: 1.827738015253393,
      maturityAPY: 1.9233043771965068,
    },
  };

  if (!parsedData.user || !parsedData.user) {
    return {};
  }
  return parsedData.user;
}

export function getBlockExplorerUrl(chainId: string, txHash?: TxHash): string {
  switch (chainId) {
    case "sifchain":
      if (!txHash) return "https://blockexplorer.sifchain.finance/";
      return `https://blockexplorer.sifchain.finance/transactions/${txHash}`;
    case "sifchain-testnet":
      if (!txHash) return `https://blockexplorer-testnet.sifchain.finance/`;
      return `https://blockexplorer-testnet.sifchain.finance/transactions/${txHash}`;
    default:
      if (!txHash) return "https://blockexplorer-devnet.sifchain.finance/";
      return `https://blockexplorer-devnet.sifchain.finance/transactions/${txHash}`;
  }
}

export function getCryptoeconomicsUrl(chainId: string): string {
  switch (chainId) {
    case "sifchain":
      return `https://api-cryptoeconomics.sifchain.finance/api`;
    case "sifchain-testnet":
      return `https://api-cryptoeconomics-devnet.sifchain.finance/api`;
    // case "sifchain-local":
    //   return `http://localhost:3000/api`; // sifnode/cryptoeconomics/js/server
    default:
      return `https://api-cryptoeconomics-devnet.sifchain.finance/api`;
  }
}

export function getRewardEarningsUrl(chainId: string): string {
  switch (chainId) {
    case "sifchain":
      return `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/netchange`;
    case "sifchain-testnet":
      return `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/netchange/testnet`;
    default:
      return `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/netchange/devnet`;
  }
}
