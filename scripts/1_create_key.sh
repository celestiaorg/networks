#!/bin/bash

# Set variables
capp=celestia-appd
CHAIN_ID="devnet-2"
NODE_NAME="eva01"
KEY_TYPE="test"

if [ $# != 1 ]; then
	echo -e "Usage:\n$0 <NODE_NAME>"
	exit 1
fi
NODE_NAME=$1

# Creating the account for validator #1
$capp keys add $NODE_NAME --keyring-backend=$KEY_TYPE
node_addr=$($capp keys show $NODE_NAME -a --keyring-backend $KEY_TYPE)

echo "----------------YOUR ACCOUNT ADDRESS BELOW----------------"
echo $node_addr

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.celestia-app/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "15s"/g' ~/.celestia-app/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.celestia-app/config/config.toml
# Open up rest api
sed -i '104 s/enable = false/enable = true/' ~/.celestia-app/config/app.toml
