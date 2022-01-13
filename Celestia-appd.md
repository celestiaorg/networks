# 1. Running a Node
Here we describe how to run and sync a node using celestia-appd.
## 1.2 Installing Dependencies
First, make sure to update and upgrade the OS:
```sh
sudo apt update && sudo apt upgrade -y
```
Install essential packages:
```sh
sudo apt install curl tar wget clang pkg-config libssl-dev jq build-essential git make ncdu -y
```
Install GO:
```sh
ver="1.17.2"
cd $HOME
wget "https://golang.org/dl/go$ver.linux-amd64.tar.gz"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "go$ver.linux-amd64.tar.gz"
rm "go$ver.linux-amd64.tar.gz"
echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> $HOME/.bash_profile
source $HOME/.bash_profile
```
To check if Go was installed correctly run:
```sh
go version
```
Output should be the version installed:
```sh
go version go1.17.2 linux/amd64
```
## 1.3 Downloading and Compiling Celestia-App
The steps below will create a binary file named celestia-appd inside `$HOME/go/bin` folder which will be used later to run the node.
```sh
cd $HOME
rm -rf celestia-app
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app/
make install
```
## 1.4 Setting up Network
First clone the networks repository:
```sh
cd $HOME
rm -rf networks
git clone https://github.com/celestiaorg/networks.git
```
To initialize the network pick a "node-name" that describes your node. The --chain-id parameter we are using here is "devnet-2". Keep in mind that this might change if a new testnet is deployed.
```sh
celestia-appd init "node-name" --chain-id devnet-2
```
Copy the `genesis.json` file. For devnet-2 we are using:
```sh
cp $HOME/networks/devnet-2/genesis.json $HOME/.celestia-app/config
```
Set seeds and peers:
```sh
SEEDS="74c0c793db07edd9b9ec17b076cea1a02dca511f@46.101.28.34:26656"
PEERS="34d4bfec8998a8fac6393a14c5ae151cf6a5762f@194.163.191.41:26656"
sed -i.bak -e "s/^seeds *=.*/seeds = \"$SEEDS\"/; s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $HOME/.celestia-app/config/config.toml
```
Reset network:
```sh
celestia-appd unsafe-reset-all
```
## 1.5 Run Ceslestia-App using Systemd
Create Celestia-App systemd file:
```sh
sudo tee <<EOF >/dev/null /etc/systemd/system/celestia-appd.service
[Unit]
Description=celestia-appd Cosmos daemon
After=network-online.target

[Service]
User=$USER
ExecStart=$HOME/go/bin/celestia-appd start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF
```
If the file was created succesfully you will be able to see its content:
```sh
cat /etc/systemd/system/celestia-appd.service
```
Enable and start celestia-appd daemon:
```sh
sudo systemctl enable celestia-appd
sudo systemctl start celestia-appd
```
Check if daemon has been started correctly:
```sh
sudo systemctl status celestia-appd
```
Check daemon logs in real time:
```sh
sudo journalctl -u celestia-appd.service -f
```
To check if your node is in sync before going forward:
```sh
curl -s localhost:26657/status | jq .result | jq .sync_info
```
Make sure that you have `"catching_up": false`, otherwise leave it running until it is in sync.

# 2. Creating a Validator
First we need to create the validator wallet. You can pick whatever wallet name you want. For our example we used "validator" as wallet name:
```sh
celestia-appd keys add validator
```
Save the mnemonic output as this is the only way to recover your validator wallet in case you lose it! For the public celestia address, fund the wallet via Discord:
```
!faucet celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
To check if tokens have arrived succesfully to the destination wallet run the command below replacing the public address with your own:
```sh
celestia-appd q bank balances celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
Create the validator on chain. Pick a MONIKER name of your choice! This is the validator name that will show up on public dashboards and explorers. The VALIDATOR_WALLET must be the same you defined previously:
```sh
MONIKER="your_moniker"
VALIDATOR_WALLET="validator"

celestia-appd tx staking create-validator \
 --amount=1000000celes \
 --pubkey=$(celestia-appd tendermint show-validator) \
 --moniker=$MONIKER \
 --chain-id=devnet-2 \
 --commission-rate=0.1 \
 --commission-max-rate=0.2 \
 --commission-max-change-rate=0.01 \
 --min-self-delegation=1000000 \
 --from=$VALIDATOR_WALLET
```
You will be prompted to confirm the transaction:
```sh
confirm transaction before signing and broadcasting [y/N]: y
```
Inputing y should provide an output similar to:
```
code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: <tx-hash>
```
You should now be able to see your validator from a block explorer such as: http://celestia.observer:3080/validators


