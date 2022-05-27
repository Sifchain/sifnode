#!/usr/bin/env bash

set -x

GENESISPATH="$HOME/.sifnoded/config/genesis.json"

echo "$(jq '.app_state.margin.params.leverage_max = "1"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.interest_rate_max = "3.000000000000000000"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.interest_rate_min = "0.005000000000000000"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.interest_rate_increase = "0.100000000000000000"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.interest_rate_decrease = "0.100000000000000000"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.health_gain_factor = "1.000000000000000000"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.epoch_length = "1"' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.pools = ["ceth", "cusdc", "cusdt", "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", "ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA", "ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D"]' $GENESISPATH)" > $GENESISPATH
echo "$(jq '.app_state.margin.params.force_close_threshold = "0.100000000000000000"' $GENESISPATH)" > $GENESISPATH
