export function cleanUpGenesisState({ remoteGenesis, defaultGenesis }) {
  return {
    ...defaultGenesis,
    ...remoteGenesis,
    initial_height: "1",
    validators: [],
    app_state: {
      ...defaultGenesis.app_state,
      ...remoteGenesis.app_state,
      auth: {
        ...defaultGenesis.app_state.auth,
        ...(remoteGenesis.app_state ? remoteGenesis.app_state.auth : {}),
        accounts: defaultGenesis.app_state.auth.accounts,
      },
      bank: {
        ...defaultGenesis.app_state.bank,
        ...(remoteGenesis.app_state ? remoteGenesis.app_state.bank : {}),
        balances: defaultGenesis.app_state.bank.balances,
        supply: defaultGenesis.app_state.bank.supply,
        denom_metadata: [],
      },
      ibc: {
        ...defaultGenesis.app_state.ibc,
        ...(remoteGenesis.app_state ? remoteGenesis.app_state.ibc : {}),
        client_genesis: defaultGenesis.app_state.ibc.client_genesis,
        connection_genesis: defaultGenesis.app_state.ibc.connection_genesis,
      },
      genutil: {
        ...defaultGenesis.app_state.genutil,
        ...(remoteGenesis.app_state ? remoteGenesis.app_state.genutil : {}),
        gen_txs: [],
      },
      ...(remoteGenesis.app_state && remoteGenesis.app_state.swap
        ? { swap: defaultGenesis.app_state.swap }
        : {}),
      ...(remoteGenesis.app_state && remoteGenesis.app_state.clp
        ? { clp: defaultGenesis.app_state.clp }
        : {}),
      ...(remoteGenesis.app_state && remoteGenesis.app_state.dispensation
        ? { dispensation: defaultGenesis.app_state.dispensation }
        : {}),
      ...(remoteGenesis.app_state && remoteGenesis.app_state.ethbridge
        ? { ethbridge: defaultGenesis.app_state.ethbridge }
        : {}),
      ...(remoteGenesis.app_state && remoteGenesis.app_state.oracle
        ? { oracle: defaultGenesis.app_state.oracle }
        : {}),
      slashing: {
        ...defaultGenesis.app_state.slashing,
        ...(remoteGenesis.app_state
          ? { params: remoteGenesis.app_state.slashing.params }
          : {}),
      },
      staking: {
        ...defaultGenesis.app_state.staking,
        ...(remoteGenesis.app_state
          ? {
              params: {
                ...remoteGenesis.app_state.staking.params,
                historical_entries: 100,
              },
            }
          : {}),
      },
      distribution: {
        ...defaultGenesis.app_state.distribution,
        ...(remoteGenesis.app_state
          ? { params: remoteGenesis.app_state.distribution.params }
          : {}),
      },
      gov: {
        ...defaultGenesis.app_state.gov,
        ...(remoteGenesis.app_state ? remoteGenesis.app_state.gov : {}),
        proposals: [],
      },
    },
  };
}
