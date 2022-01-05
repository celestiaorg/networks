#!/bin/bash

set -e
if [ "$1" == "--help" ] || [ -z "$2" ]; then
    echo "Usage: ./fix_genesis.sh path/to/genesis.json new_token_name"
    exit 0
fi

# $1 -> path to genesis.json
# $2 -> becomes an token arg in jq 
cat $1 | jq --arg token $2 '.app_state.crisis.constant_fee.denom = $token' | sponge $1
cat $1 | jq --arg token $2 '.app_state.mint.params.mint_denom = $token' | sponge $1
cat $1 | jq --arg token $2 '.app_state.staking.params.bond_denom = $token' | sponge $1
