# Run a Celestia Bridge Node

- [Run a Celestia Bridge Node](#run-a-celestia-bridge-node)
  - [Dependencies](#dependencies)
    - [Update Packages](#update-packages)
  - [Installing GO](#installing-go)
  - [Part 1: Deploy the Celestia App](#part-1-deploy-the-celestia-app)
    - [Install Celestia App](#install-celestia-app)
    - [Set up P2P Network](#set-up-p2p-network)
    - [Run Celestia-App using Systemd](#run-celestia-app-using-systemd)
    - [Create a Wallet](#create-a-wallet)
    - [Fund a Wallet](#fund-a-wallet)
    - [Delegate Stake to a Validator](#delegate-stake-to-a-validator)
  - [Part 2: Deploy the Celestia Node](#part-2-deploy-the-celestia-node)
    - [Install Celestia Node](#install-celestia-node)
    - [Get the trusted hash](#get-the-trusted-hash)
    - [Initialize the Bridge Node](#initialize-the-bridge-node)
    - [Configure the Bridge Node](#configure-the-bridge-node)
    - [Start the Bridge Node](#start-the-bridge-node)
  - [Run a Validator Bridge Node](#run-a-validator-bridge-node)

## Dependencies

### Update Packages
First, make sure to update and upgrade the OS:
```sh
sudo apt update && sudo apt upgrade -y
```
These are essential packages which are necessary execute many tasks like downloading files, compiling and monitoring the node:
```sh
sudo apt install curl tar wget clang pkg-config libssl-dev jq build-essential git make ncdu -y
```
## Installing GO
It is necessary to install the GO language in the OS, so we can later compile the Celestia Application. On our example, we are using version 1.17.2:
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

## Part 1: Deploy the Celestia App
This section describes part 1 of Celestia Bridge Node setup: running a Celestia App daemon with an internal Celestia Core node.

> Caveat: Make sure you have at least 100+ Gb of free space to safely install+run the Bridge Node.

### Install Celestia App
The steps below will create a binary file named celestia-appd inside `$HOME/go/bin` folder which will be used later to run the node.
```sh
cd $HOME
rm -rf celestia-app
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app/
git checkout tags/v0.1.0 -b v0.1.0
make install
```
To check if the binary was successfully compiled you can run the binary using the `--help` flag:
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
### Set up P2P Network
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
### Run Celestia-App using Systemd
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
If the file was created successfully you will be able to see its content:
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

### Create a Wallet
You can pick whatever wallet name you want. For our example we used "validator" as the wallet name:
```sh
celestia-appd keys add validator
```
Save the mnemonic output as this is the only way to recover your validator wallet in case you lose it! 

### Fund a Wallet
For the public celestia address, you can fund the previously created wallet via Discord by sending this message to #faucet channel:
```
!faucet celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
Wait to see if you get a confirmation that the tokens have been successfully sent. To check if tokens have arrived successfully to the destination wallet run the command below replacing the public address with your own:
```sh
celestia-appd q bank balances celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Delegate Stake to a Validator
If you want to delegate more stake to any validator, including your own you will need the `celesvaloper` address of the validator in question. You can either check it using the block explorer mentioned above or you can run the command below to get the `celesvaloper` of your local validator wallet in case you want to delegate more to it:
```sh
celestia-appd keys show $VALIDATOR_WALLET --bech val -a
```
After entering the wallet passphrase you should see a similar output:
```sh
Enter keyring passphrase:
celesvaloper1q3v5cugc8cdpud87u4zwy0a74uxkk6u43cv6hd
```
To delegate tokens to the the `celesvaloper` validator, as an example you can run:
```sh
celestia-appd tx staking delegate celesvaloper1q3v5cugc8cdpud87u4zwy0a74uxkk6u43cv6hd 1000000celes --from=$VALIDATOR_WALLET --chain-id=devnet-2
```
If successful, you should see a similar output as:
```sh
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
You can check if the TX hash went through using the block explorer by inputting the `txhash` ID that was returned.

## Part 2: Deploy the Celestia Node
This section describes part 2 of Celestia Bridge Node setup: running a Celestia Node daemon.

### Install Celestia Node
Install the Celestia Node binary, which will be used to run the Bridge Node.
```sh
cd $HOME
rm -rf celestia-node
git clone https://github.com/celestiaorg/celestia-node.git
cd celestia-node/
make install
```
Verify that the binary is working and check the version with `celestia version` command:

```console
$ celestia version
Semantic version: v0.2.0
Commit: 1fcf0c0bb5d5a4e18b51cf12440ce86a84cf7a72
Build Date: Fri 04 Mar 2022 01:15:07 AM CET
System version: amd64/linux
Golang version: go1.17.5
```

### Get the trusted hash
> Caveat: You need a running celestia-app in order to continue this guideline. Please refer to [celestia-app.md](https://github.com/celestiaorg/networks/celestia-app.md) for installation.

You need to have the trusted hash in order to initialize the Celestia Bridge Node. In order to know the hash, you need to query the Celestia App:
```sh
curl -s http://localhost:26657/block?height=1 | grep -A1 block_id | grep hash
```

### Initialize the Bridge Node
```sh
celestia bridge init --core.remote <ip:port of celestia-app> --headers.trusted-hash <hash_from_celestia_app>
```

Example:
```sh 
celestia bridge init --core.remote tcp://127.0.0.1:26657 --headers.trusted-hash 4632277C441CA6155C4374AC56048CF4CFE3CBB2476E07A548644435980D5E17
```

### Configure the Bridge Node

In order for your Celestia Bridge Node to communicate with other Bridge Ndoes, then you need to add them as `mutual peers` in the `config.toml` file and allow the peer exchange. Please navigate to `networks/devnet-2/celestia-node/mutual_peers.txt` to find the list of mutual peers

For more information on `config.toml`, please navigate to [this link](./config-toml.md)
```sh
nano ~/.celestia-bridge/config.toml
```
```sh
...
[P2P]
  ...
  #add multiaddresses of other celestia bridge nodes
  
  MutualPeers = [
    "/ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22", 
    "/ip4/x.x.x.x/tcp/yyy/p2p/abc"] #the /ip4/x.x.x.x is only for example. Don't add it! 
  PeerExchange = true #change this line to true. Be default it's false
  ...
...
```

### Start the Bridge Node
```sh
celestia bridge start
```
Now, the Celestia bridge node will start syncing headers and storing blocks from Celestia application. 

> Note: At startup, we can see the `multiaddress` from Celestia Bridge Node. This is <b>needed for future Light Node</b> connections and communication between Celestia Bridge Nodes

Example:
```sh
/ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22
```

---

## Run a Validator Bridge Node
Optionally, if you want to join the active validator list, you can create your own validator on-chain following the instructions below. Keep in mind that these steps are necessary ONLY if you want to participate in the consensus.

Pick a MONIKER name of your choice! This is the validator name that will show up on public dashboards and explorers. VALIDATOR_WALLET must be the same you defined previously. Parameter `--min-self-delegation=1000000` defines the amount of tokens that are self delegated from your validator wallet.
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
Inputting y should provide an output similar to:
```sh
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
You should now be able to see your validator from a block explorer such as: https://celestia.observer/validators

If you want to run a Celestia Node, check the documentation [here](./celestia-node.md).
