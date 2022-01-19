- [Running a Non-Validating Celestia Application Node](#running-a-non-validating-celestia-application-node)
  - [Installing Dependencies](#installing-dependencies)
  - [Installing GO](#installing-go)
  - [Downloading and Compiling Celestia-App](#downloading-and-compiling-celestia-app)
  - [Setting up Network](#setting-up-network)
  - [Running Ceslestia-App using Systemd](#running-ceslestia-app-using-systemd)
- [Running a Validator](#running-a-validator)
  - [Creating a Validator Wallet](#creating-a-validator-wallet)
  - [Funding the Validator Wallet](#funding-the-validator-wallet)
  - [Creating the Validator On-Chain](#creating-the-validator-on-chain)

# Running a Non-Validating Celestia Application Node
This chapter describes how to run a Non-Validating Celestia Application Node. Non-Validating nodes allow you to interact with the Celestia Network without having to create and maintain a validator.

## Installing Dependencies
First, make sure to update and upgrade the OS:
```sh
sudo apt update && sudo apt upgrade -y
```
These are essential packages which are necessary execute many tasks like downloading files, compiling and monitoring the node:
```sh
sudo apt install curl tar wget clang pkg-config libssl-dev jq build-essential git make ncdu -y
```
## Installing GO
It is necessary to install the GO language in the OS so we can later compile the Celestia Application. On our example, we are using version 1.17.2:
```sh
ver="1.17.2"
cd $HOME
wget "https://golang.org/dl/go$ver.linux-amd64.tar.gz"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "go$ver.linux-amd64.tar.gz"
rm "go$ver.linux-amd64.tar.gz"
```
Now we need to add the `/usr/local/go/bin` directory to `$PATH`:
```
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
## Downloading and Compiling Celestia-App
The steps below will create a binary file named celestia-appd inside `$HOME/go/bin` folder which will be used later to run the node.
```sh
cd $HOME
rm -rf celestia-app
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app/
make install
```
To check if the binary was succesfully compiled you can run the binary using the `--help` flag:
```sh
cd $HOME/go/bin
./celestia-appd --help
```
You should see a similar output:
```
Stargate CosmosHub App

Usage:
  celestia-appd [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  config              Create or query an application CLI configuration file
  debug               Tool for helping with debugging your application
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, and application configuration files
  keys                Manage your application's keys
  migrate             Migrate genesis to a specified target version
  query               Querying subcommands
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  tx                  Transactions subcommands
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information

Flags:
  -h, --help                help for celestia-appd
      --home string         directory for config and data (default "/home/pops/.celestia-app")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors

Use "celestia-appd [command] --help" for more information about a command.
```
## Setting up Network
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
## Running Ceslestia-App using Systemd
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

# Running a Validator
Optionally, if you want to join the active validator list, you can create your own validator on-chain following the instructions below. Keep in mind that these steps are necessary ONLY if you want to participate in the consensus.

## Creating a Validator Wallet
First we need to create the validator wallet. You can pick whatever wallet name you want. For our example we used "validator" as wallet name:
```sh
celestia-appd keys add validator
```
Save the mnemonic output as this is the only way to recover your validator wallet in case you lose it! 
## Funding the Validator Wallet
For the public celestia address, you can fun the validator wallet via Discord sending this message to #faucet channel:
```
!faucet celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
Wait to see if you get a confirmation that the tokens have been successfully sent. To check if tokens have arrived succesfully to the destination wallet run the command below replacing the public address with your own:
```sh
celestia-appd q bank balances celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
## Creating the Validator On-Chain
Pick a MONIKER name of your choice! This is the validator name that will show up on public dashboards and explorers. The VALIDATOR_WALLET must be the same you defined previously. Parameter `--min-self-delegation=1000000` defines the amount of tokens that are self delegated from your validator wallet.
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


