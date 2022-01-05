Run Celestia Devnet
---
- [Run Celestia Devnet](#run-celestia-devnet)
- [Pre-Requisites](#pre-requisites)
  - [Installed Celestia App and Celestia Node](#installed-celestia-app-and-celestia-node)
  - [Fork/Clone networks repo](#forkclone-networks-repo)
- [Running a Celestia App](#running-a-celestia-app)
  - [Legend:](#legend)
    - [- Facilitator (F) - orchestrates genesis setup](#--facilitator-f---orchestrates-genesis-setup)
    - [- Others (O) - other participants](#--others-o---other-participants)
  - [Phase 0: Environment ideas/hints](#phase-0-environment-ideashints)
  - [Phase 1: Creating accounts and initializing celestia app chain](#phase-1-creating-accounts-and-initializing-celestia-app-chain)
    - [(O) steps:](#o-steps)
    - [(F) Steps:](#f-steps)
  - [Phase 2: GenTx](#phase-2-gentx)
    - [(O) steps:](#o-steps-1)
    - [(F) Steps:](#f-steps-1)
  - [Phase 3: Launch](#phase-3-launch)
    - [(O) steps:](#o-steps-2)
    - [Everybody:](#everybody)
  - [Verifying Celestia App Network](#verifying-celestia-app-network)
  - [Phase X: Be a Validator in the running `devnet-2` network:](#phase-x-be-a-validator-in-the-running-devnet-2-network)
  - [Troubleshooting Validator Node setup](#troubleshooting-validator-node-setup)
- [Running a Celestia Node](#running-a-celestia-node)
  - [Pre-Requisites](#pre-requisites-1)
  - [Full Node Configuration](#full-node-configuration)
  - [Light Client Configuration](#light-client-configuration)
- [Data Availability Sampling(DAS)](#data-availability-samplingdas)
  - [Pre-Requisites:](#pre-requisites-2)
  - [Legend:](#legend-1)
    - [- First terminal(FT) run light client with logs info](#--first-terminalft-run-light-client-with-logs-info)
    - [- Second terminal(ST) submit payForMessage tx](#--second-terminalst-submit-payformessage-tx)
  - [Steps:](#steps)
## Pre-Requisites 
### Installed Celestia App and Celestia Node
If you haven't installed either of app/node, please navigate to each of them
- [celestia-app repo](https://github.com/celestiaorg/celestia-app)
- [celestia-node repo](https://github.com/celestiaorg/celestia-node)

### Fork/Clone networks repo
```sh
git clone https://github.com/<your_github>/networks.git
```

## Running a Celestia App
### Legend: 
#### - Facilitator (F) - orchestrates genesis setup
#### - Others (O) - other participants

### Phase 0: Environment ideas/hints
The following instructions assume that you are using Digital Ocean as a cloud provider. Though it should be straightforward to use another provider or even deploy on a local bare-metal server instead. 
Tested setups:
1. One Digital Ocean Droplet on Ubuntu(1 CPU and 1Gb should be enough) and your local Linux-based(Ubuntu prefferably) OS
2. Two Digital Ocean Droplets on Ubuntu(1 CPU and 1 Gb should be enough)
3. Remember to update your firewall settings. You can check and execute `.networks/scripts/firewall_ubuntu.sh` script

<u>Please don't try to build `celestia app` on macOS Big Sur and higher due to this issue [#134](https://github.com/celestiaorg/celestia-app/issues/134)</u>


### Phase 1: Creating accounts and initializing celestia app chain
#### (O) steps: 
1. Creates an account using a script in networks repository
```sh
cd networks/scripts
node_name="MightyValidator"
./1_create_key.sh $node_name
```
2. After execution of the script, (O) passes account's address to (F)

#### (F) Steps:
1. Initializes the celestia app chain using a script in networks repository
```sh
cd networks/scripts
./creator.sh
```
2. Adds (O)'s account address using command
```sh
celestia-appd add-genesis-account $O_acc_addr 800000000000celes
```
<b><u>Note: Currently, we've found an issue that `genesis.json` is using `stake` instead of `celes` token [#158](https://github.com/celestiaorg/celestia-app/issues/158)</u></b>
<b>Please, either manually change `stake` to `celes` in `genesis.json` or execute this script:</b>
```sh
sudo apt-get install moreutils jq #contains sponge that we need
cd networks/scripts
./fix_genesis.sh ~/.celestia-app/config/genesis.json
```


3. Send the `genesis.json` to (O). Location of genesis.json file can be find by these default path:
```sh
cd ~/.celestia-app/config/genesis.json
```

### Phase 2: GenTx
#### (O) steps: 
1. Copies `genesis.json` from (F) to local `.celestia-app/config/` directory:
```sh
cp <path_to_downloaded_genesis.json>/genesis.json  ~/.celestia-app/config/
```
2. Executes `gentx` command in CLI:
```sh
celestia-appd gentx $node_name 5000000000celes --keyring-backend=test --chain-id devnet-2
``` 
3. Passes `gentx-<hash>.json` file to (F). Location of the file can be find in this default path: 
```sh
cd .celestia-app/config/gentx/ && ls
```
#### (F) Steps:
1. Copies `gentx-<hash>.json` from (O) to local `.celestia-app/config/gentx` directory
2. Executes command in CLI:
```sh
celestia-appd collect-gentxs
```
3. Send the updated version of `genesis.json` to (O)

### Phase 3: Launch
#### (O) steps:
1. Copies <b><i>updated</i></b> `genesis.json` from (F) to local `.celestia-app/config/` directory:
```sh
cp <path_to_downloaded_genesis.json>/genesis.json  ~/.celestia-app/config/
```
#### Everybody:  
```sh
celestia-appd start
```
> :grey_exclamation:	 If your node is having trouble finding peers to connect with, please update your [seeds](devnet-2/seeds.txt) and [peers](devnet-2/peers.txt) in your config.toml

### Verifying Celestia App Network

To check that all validators are communicating with each other correctly, please look for a same block height for each of a running validator (let’s say block height = 100) 

To confirm that everything is ok, you have to check that commit hashes for same block height are the same for all validators.

### Phase X: Be a Validator in the running `devnet-2` network:
1. Creates an account using a script in networks repository
```sh
cd networks/scripts
node_name="chani"
./1_create_key.sh $node_name
```
2. Request some testnet tokens in [Discord Faucet Channel](https://discord.gg/NEfhSscRSx)
3. Executes command in CLI:

```sh
celestia-appd init <moniker_name> --chain-id devnet-2
```
Example: 
```sh 
celestia-appd init fremen --chain-id devnet-2
```
4. Copies `genesis.json` from `networks/devnet-2` repo to local `.celestia-app/config/` directory:
```sh
cp networks/devnet-2/genesis.json  ~/.celestia-app/config/
```
5. Starts syncing the chain using command in CLI:
```sh 
celestia-appd start --p2p.seeds 74c0c793db07edd9b9ec17b076cea1a02dca511f@46.101.28.34:26656
```
Syncing finishes around 1-2 hours

6. Execute command:
```sh
celestia-appd tx staking create-validator \
 --amount=1000000000celes \
 --pubkey=$(celestia-appd tendermint show-validator) \
 --moniker=space \
 --chain-id=devnet-2 \
 --commission-rate=0.1 \
 --commission-max-rate=0.2 \
 --commission-max-change-rate=0.01 \
 --min-self-delegation=1000000000 \
 --from=$node_name \
 --keyring-backend=test
```

which should prompt for:

```sh
confirm transaction before signing and broadcasting [y/N]:
```

Inputting `y` should provide an output similar to:

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

You should now be able to see your validator from a block explorer such as http://celestia.observer:3080/validators.

### Troubleshooting Validator Node setup

If you get an error such as 

```Error: <keyname>: key not found```,

this means your key, the field referenced by the `--from` option, does not exist.

You can fix this by adding your key manually to the keyring via:

```
celestia-appd keys add --recover <keyname>
```

followed by a prompt to enter a bip39 mnemonic, which is the mnemonic that was created as part of `1_create_key.sh` script in the first step.

You'll also be asked for a passphrase which is an input you have to define.

After this fix repeat step 6.


## Running a Celestia Node
### Pre-Requisites
You need to have the trusted hash in order to initialize the full celestia node
In order to know the hash, you need to query the celestia-app:
```sh
curl -s http://localhost:26657/block?height=1 | grep -A1 block_id | grep hash
```

### Full Node Configuration
1. Initialize the full node
```sh
celestia full init --core.remote <ip:port of celestia-app> --headers.trusted-hash <hash_from_celestia_app>
```

Example:
```sh 
celestia full init --core.remote tcp://127.0.0.1:26657 --headers.trusted-hash 3BDFBD7E2D97215CFA600DD8B39AAEECC015E43FEE7B8A4D8A8B630E8B4D4006
```

2. Edit configurations

If you have multiple celestia full nodes, then you need to add them as mutual peers in the `config.toml` file and allow the peer exchange. This is needed in order for celestia full nodes communication between each other.
```sh
nano ~/.celestia-full/config.toml
```
```sh
...
[P2P]
  ...
  MutualPeers = ["/ip4/<ip>/tcp/2121/p2p/<pubKey>"] #add multiaddresses of other celestia full nodes
  PeerExchange = true #change this line to true. Be default it's false
  ...
...
```

3. Start the full node
```sh
celestia full start
```
Now, the celestia full node will start syncing headers and storing blocks from celestia app. 

<u>Note: At startup, we can see the `multiaddress` from celestia full node. This is <b>needed for future light client</b> connections and communication between celestia full nodes</u>

Example:
```sh
/ip4/46.101.245.50/tcp/2122/p2p/12D3KooWKNBZvF93L92aTYs6jRyozicGmuu3cF9gotMtUCeHAPYN
```

### Light Client Configuration
To start the light client, we need to know 2 variables:
- Trusted peer’s multi address to connect to (a celestia full node is the case here)
- Trusted block hash from celestia-app

1. Initialize the light client

```sh
celestia full init --headers.trusted-peer <full_node_multiaddress> --headers.trusted-hash <hash_from_celestia_app>
```

Example: 

```sh 
celestia light init --headers.trusted-peer /ip4/46.101.245.50/tcp/2122/p2p12D3KooWKNBZvF93L92aTYs6jRyozicGmuu3cF9gotMtUCeHAPYN --headers.genesis-hash F10679041E3D55363405C1D3080B91004BCAE471F35F1FBABC345B8237AEFDA2 
```

2. Start the light client
```sh
celestia light start
```
Now, the celestia light client will start syncing headers and do data availability sampling(DAS) from the full node.

## Data Availability Sampling(DAS)

### Pre-Requisites:
This is a list of runnining components you need in order to successfully continue this chapter:
- celestia app validator
- celestia full node
- celestia light client

### Legend:
You will need 2 terminals in order to see how DASing works:
#### - First terminal(FT) run light client with logs info
#### - Second terminal(ST) submit payForMessage tx

### Steps:
1. In (ST) Submit a `payForMessage` transaction with `celestia-appd`
```sh
celestia-appd tx payment payForMessage <hex_namespace> <hex_message> --from <node_name> --keyring-backend <keyring-name> --chain-id <chain_name> -y
```
Example:
```sh 
celestia-appd tx payment payForMessage 0102030405060708 68656c6c6f43656c6573746961444153 --from eva00 --keyring-backend test --chain-id devnet-2 -y
```
2. In (FT) you should see in logs how DAS is working

Example:
```sh
INFO	das	das/daser.go:96	sampling successful	{"height": 81547, "hash": "DE0B0EB63193FC34225BD55CCD3841C701BE841F29523C428CE3685F72246D94", "square width": 2, "finished (s)": 0.000117466}
```
