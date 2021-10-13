export async function startRelayers({ chainsProps }) {
  const { sifchain: sifChainProps, ...otherChainsProps } = chainsProps;

  return Promise.all(
    Object.values(otherChainsProps).map(async ({ home }) => {
      const relayerHome = `${home}/relayer`;

      const proc = await nothrow(
        $`ibc-relayer start -v --poll 10 --home ${relayerHome}`
      );

      return {
        proc,
      };
    })
  );
}
