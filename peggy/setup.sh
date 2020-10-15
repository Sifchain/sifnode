#
rm -rf ~/.sifnoded 
rm -rf ~/.sifnodecli

sifnoded init local --chain-id=peggy

sifnodecli config keyring-backend test

sifnodecli config chain-id peggy
sifnodecli config trust-node true
sifnodecli config indent true
sifnodecli config output json

# Create a key to hold your validator account and for another test account
sifnodecli keys add validator
sifnodecli keys add testuser

# Initialize the genesis account and transaction
sifnoded add-genesis-account $(sifnodecli keys show validator -a) 1000000000stake,1000000000atom
sifnodecli tx send validator $(sifnodecli keys show testuser -a) 10atom --yes

# Create genesis transaction
sifnoded gentx --name validator --keyring-backend test

# Collect genesis transaction
sifnoded collect-gentxs


