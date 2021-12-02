Run celestia-app chain
---

## Pre-Requisites 
### Clone celestia-app repo and install
```sh
git clone https://github.com/celestiaorg/celestia-app.git

cd celestia-app 

make install

# check that celestia-appd --help is working properly
celestia-appd --help
```
### Fork/Clone networks repo
```sh
git clone https://github.com/<your_github>/networks.git
```

## Legend: 
### - Facilitator (F) - orchestrates genesis setup
### - Others (O) - other participants

## Phase 0: Environment ideas/hints
If you want to run celestia app chain, please consider using such setups: 

1. One Digital Ocean Droplet on Ubuntu(1 CPU and 1Gb should be enough) and your local Linux-based(Ubuntu prefferably) OS
2. Two Digital Ocean Droplets on Ubuntu(1 CPU and 1 Gb should be enough)
3. Docker Ubuntu based images
4. Remember to update your firewall settings. You can check and execute `.networks/scripts/firewall_ubuntu.sh` script

<u>Please don't try to build `celestia app` on macOS Big Sur and higher due to this issue [#134](https://github.com/celestiaorg/celestia-app/issues/134)</u>


## Phase 1: Creating accounts and initializing celestia app chain
### (O) steps: 
1. Creates an account using a script in networks repository
```sh
cd networks/scripts
node_name="MightyValidator"
./1_create_key.sh $node_name
```
2. After execution of the script, (O) passes account's address to (F)

### (F) Steps:
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
sudo apt-get install moreutils #contains sponge that we need
cd networks/scripts
./fix_genesis.json
```


3. Send the `genesis.json` to (O). Location of genesis.json file can be find by these default path:
```sh
cd ~/.celestia-app/config/genesis.json
```

## Phase 2: GenTx
### (O) steps: 
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
### (F) Steps:
1. Copies `gentx-<hash>.json` from (O) to local `.celestia-app/config/gentx` directory
2. Executes command in CLI:
```sh
celestia-appd collect-gentxs
```
3. Send the updated version of `genesis.json` to (O)

## Phase 3: Launch
### (O) steps:
1. Copies <b><i>updated</i></b> `genesis.json` from (F) to local `.celestia-app/config/` directory:
```sh
cp <path_to_downloaded_genesis.json>/genesis.json  ~/.celestia-app/config/
```
### Everybody:  
```sh
celestia-appd start
```

## Verifying Celestia App Network

To check that all validators are communicating with each other correctly, please look for a same block height for each of a running validator (let’s say block height = 100) 

To confirm that everything is ok, you have to check that commit hashes for same block height are the same for all validators.
