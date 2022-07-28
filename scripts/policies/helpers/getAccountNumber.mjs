import { spinner } from "zx/experimental";
import { binary, queryFlags, keyringFlags } from "./constants.mjs";

export async function getAccountNumber(key) {
  const { account_number: accountNumber, sequence } = await spinner(
    "retrieving account number and sequence                              ",
    () =>
      within(async () => {
        // $.verbose = true;
        const response = await $`\
${binary} q auth account \
  $(${binary} keys show ${key} -a ${keyringFlags}) \
  ${queryFlags} \
  --output=json \
  | jq -r '{account_number: .account_number, sequence: .sequence}'`;
        return JSON.parse(response);
      })
  );
  return { accountNumber, sequence };
}
