const { SIFNODE_CHAIN_ID } = process.env;

export function checkEnv() {
  if (SIFNODE_CHAIN_ID !== "localnet") {
    throw new Error(
      `wrong environment file use localnet instead of ${SIFNODE_CHAIN_ID}`
    );
  }
}
