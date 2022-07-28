import { spinner } from "zx/experimental";

export async function getTokens(entries) {
  const tokens = await spinner(
    "retrieving tokens                              ",
    () =>
      within(() =>
        entries
          .filter(
            ({ denom }) =>
              denom.startsWith("c") ||
              denom.startsWith("ibc/") ||
              denom === "rowan"
          )
          .map(({ denom }) => denom)
      )
  );
  return tokens;
}
