#!/usr/bin/env bash

dir=~/.celestia-app/config/gentx

if [ "$1" != "" ]; then
	dir=$1
fi

echo "Adding genesis accounts from files located in: $dir"
for f in `ls $dir/gentx*.json`;
do 
	echo $f
	addr=`cat $f | jq '.body.messages[] | .delegator_address' | tr -d \"`
	amount=`cat $f | jq '.body.messages[] | .value.amount' | tr -d \"`
	denom=`cat $f | jq '.body.messages[] | .value.denom' | tr -d \"`
	celestia-appd add-genesis-account $addr "$amount$denom"
done
