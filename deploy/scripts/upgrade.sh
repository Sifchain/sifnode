#!/usr/bin/env sh
#
# Sifnode migration (from Cosmos 0.39.x to 0.42.x).
#

#
# Usage.
#
usage() {
  cat <<- EOF
  Usage: $0 [OPTIONS]

  Options:
  -h      This help output.
  -b      Block height to export the current state from.
  -c      New Chain ID.
  -s      Cosmos SDK target version.
  -t      Genesis time (in UTC).
  -v      The new sifnoded binary version.

EOF
  exit 1
}

#
# Setup
#
setup() {
  set_block_height "${1}"
  set_chain_id "${2}"
  set_cosmos_sdk_version "${3}"
  set_genesis_time "${4}"
  set_version "${5}"
  create_export_state_dir
}

#
# Already upgraded?
#
upgraded() {
  if [ -f "${HOME}"/.sifnoded/."${COSMOS_SDK_VERSION}"_upgraded ]; then
    exit 0
  fi
}

#
# Set block height.
#
set_block_height() {
  BLOCK_HEIGHT=${1}
}

#
# Set Chain ID.
#
set_chain_id() {
  CHAIN_ID=${1}
}

#
# Set Genesis time.
#
set_cosmos_sdk_version() {
  COSMOS_SDK_VERSION=${1}
}

#
# Set Genesis time.
#
set_genesis_time() {
  GENESIS_TIME=${1}
}

#
# Set version.
#
set_version() {
  VERSION=${1}
}

#
# Create export state dir.
#
create_export_state_dir() {
  EXPORT_STATE_DIR="${HOME}"/.sifnoded/"${COSMOS_SDK_VERSION}"_exports
  mkdir "${EXPORT_STATE_DIR}"
}

#
# Backup.
#
backup() {
  BACKUP_DIR="${HOME}"/.sifnoded/backups/"${BLOCK_HEIGHT}"/
  mkdir -p "${BACKUP_DIR}"
  cp -avr "${HOME}"/.sifnoded/data/ "${BACKUP_DIR}"
  cp -avr "${HOME}"/.sifnoded/config/ "${BACKUP_DIR}"
}

#
# Export state.
#
export_state() {
  "${HOME}"/.sifnoded/cosmovisor/current/bin/sifnoded export --for-zero-height --height "${BLOCK_HEIGHT}" > "${EXPORT_STATE_DIR}"/exported_state.json
}

#
# Migrate exported state.
#
migrate_exported_state() {
  # Need to be the latest binary.
  "${HOME}"/.sifnoded/cosmovisor/"${VERSION}"/bin/sifnoded migrate v"${COSMOS_SDK_VERSION}" "${EXPORT_STATE_DIR}"/exported_state.json \
    --chain-id "${CHAIN_ID}" \
    --genesis-time "${GENESIS_TIME}" > "${EXPORT_STATE_DIR}"/migrated_state.json
}

#
# Configure IBC
#
configure_ibc() {
  cat "${EXPORT_STATE_DIR}"/migrated_state.json | jq '.app_state |= . + {"ibc":{"client_genesis":{"clients":[],"clients_consensus":[],"create_localhost":false},"connection_genesis":{"connections":[],"client_connection_paths":[]},"channel_genesis":{"channels":[],"acknowledgements":[],"commitments":[],"receipts":[],"send_sequences":[],"recv_sequences":[],"ack_sequences":[]}},"transfer":{"port_id":"transfer","denom_traces":[],"params":{"send_enabled":false,"receive_enabled":false}},"capability":{"index":"1","owners":[]}}' > "${EXPORT_STATE_DIR}"/genesis_ibc.json
  mv "${EXPORT_STATE_DIR}"/genesis_ibc.json "${EXPORT_STATE_DIR}"/genesis.json
}

#
# Reset old state.
#
reset_old_state() {
  "${HOME}"/.sifnoded/cosmovisor/"${VERSION}"/bin/sifnoded unsafe-reset-all
}

#
# Install genesis.
#
install_genesis() {
  cp "${EXPORT_STATE_DIR}"/genesis.json "${HOME}"/.sifnoded/config/genesis.json
}

#
# Update config.
#
update_config() {
  printf '\n[api]\nenable = false\nswagger = false\n\n[grpc]\nenable = false\naddress = "0.0.0.0:9090"\n\n[state-sync]\nsnapshot-interval = 0\nsnapshot-keep-recent = 2\n\n' >> "${HOME}"/.sifnoded/config/config.toml
}

#
# Update symlink
#
update_symlink() {
  rm "${HOME}"/.sifnoded/cosmovisor/current
  ln -s "${HOME}"/.sifnoded/cosmovisor/upgrades/"${VERSION}" "${HOME}"/.sifnoded/cosmovisor/current
}

#
# Completed.
#
completed() {
  touch "${HOME}"/.sifnoded/."${COSMOS_SDK_VERSION}"_upgraded
}

#
# Run.
#
run() {
  # Setup.
  printf "\nConfiguring environment for upgrade..."
  setup "${1}" "${2}" "${3}" "${4}" "${5}"

  # Backup.
  printf "\nTaking a backup..."
  backup

  # Check if already upgraded?
  printf "\nChecking if validator has already been upgraded..."
  upgraded

  # Export state.
  printf "\nExporting the current state..."
  export_state

  # Migrate exported state.
  printf "\nMigrating the exported state..."
  migrate_exported_state

  # Configure IBC.
  printf "\nConfiguring IBC..."
  configure_ibc

  # Reset old state.
  printf "\nResetting old state..."
  reset_old_state

  # Install the new genesis.
  printf "\nInstalling the new genesis file..."
  install_genesis

  # Updating the config.
  printf "\nUpdating the node config (api,grpc,state-sync)..."
  update_config

  # Update symlink.
  printf "\nUpdating the cosmovisor symlink..."
  update_symlink

  # Complete.
  printf "\nUpgrade complete! Good luck!"
  completed
}

# Check the supplied opts.
while getopts ":hb:c:s:t:v:" o; do
  case "${o}" in
    h)
      usage
      ;;
    b)
      b=${OPTARG}
      ;;
    c)
      c=${OPTARG}
      ;;
    s)
      s=${OPTARG}
      ;;
    t)
      t=${OPTARG}
      ;;
    v)
      v=${OPTARG}
      ;;
    *)
      usage
      ;;
  esac
done
shift $((OPTIND-1))

if [ -z "${b}" ]; then
  usage
fi

if [ -z "${c}" ]; then
  usage
fi

if [ -z "${s}" ]; then
  usage
fi

if [ -z "${t}" ]; then
  usage
fi

if [ -z "${v}" ]; then
  usage
fi

# Run.
run "${b}" "${c}" "${s}" "${t}" "${v}"
