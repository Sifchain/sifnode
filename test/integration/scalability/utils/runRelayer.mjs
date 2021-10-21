export async function runRelayer({ home }) {
  const relayerHome = `${home}/relayer`;

  const proc = await nothrow(
    $`ibc-relayer start -v --poll 10 --home ${relayerHome}`
  );

  return {
    proc,
  };
}
