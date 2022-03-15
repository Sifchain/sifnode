# usage:

# $0 1000 /tmp/somedirectory

# Creates 1000 sif accounts in a test keychain located in /tmp/somedirectory

naccounts=$1
shift
home=$1
shift

for i in $(seq $naccounts)
do
  sifnoded keys add $(uuidgen) --keyring-backend test --home $home
done
