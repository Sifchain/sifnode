export function getDenoms() {
  return {
    entries: [
      //   {
      //     is_whitelisted: true,
      //     denom: "ceth",
      //     decimals: 18,
      //     ibc_counterparty_denom: "xeth",
      //     permissions: ["CLP", "IBCEXPORT", "IBCIMPORT"],
      //   },
      //   {
      //     is_whitelisted: true,
      //     denom: "xeth",
      //     decimals: 10,
      //     unit_denom: "ceth",
      //     permissions: ["CLP", "IBCEXPORT", "IBCIMPORT"],
      //   },
      //   {
      //     is_whitelisted: true,
      //     denom: "cdash",
      //     decimals: 18,
      //     permissions: ["CLP", "IBCEXPORT", "IBCIMPORT"],
      //   },
      {
        is_whitelisted: true,
        denom: "rowan",
        decimals: 18,
        ibc_counterparty_denom: "xrowan",
        permissions: ["CLP", "IBCEXPORT", "IBCIMPORT"],
      },
      {
        is_whitelisted: true,
        denom: "xrowan",
        decimals: 10,
        unit_denom: "rowan",
        permissions: ["CLP", "IBCEXPORT", "IBCIMPORT"],
      },
    ],
  };
}
