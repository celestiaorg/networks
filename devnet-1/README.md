# Run celestia-app & core

## clone the repo and install
```sh
git clone https://github.com/celestiaorg/lazyledger-app.git

cd lazyledger-app 

make install

# check that celestia-appd --help is working properly
celestia-appd --help
```

## Initialize your validator

```sh
celestia-appd config chain-id devnet-1

# provide network and moniker name (i.e. name of your validator -- choose any name!)
celestia-appd init <moniker_name> --chain-id devnet-1
```
This will create a new `.celestia-app` folder in your HOME directory.

### Download Pregenesis File

You can now download the "pregenesis" file for the chain.  This is a genesis file with the chain-id and airdrop balances.

```sh
cd $HOME/.celestia-app/config/

curl https://raw.githubusercontent.com/celestiaorg/networks/master/devnet-1/pregenesis.json > $HOME/.celestia-app/config/genesis.json
```


## Get your validator pubkey 

You must get your validator's consensus pubkey as it will be necessary to include in the transaction to create your validator.

If you are using Tendermint's native `priv_validator.json` as your consensus key, you display your validator public key using the following command

```
celestia-appd tendermint show-validator
```

The pubkey should be formatted with the bech32 prefix `celesvalconspub1`.

If you are using a custom signing mechanism such as `tmkms`, please refer to their relevant docs to retrieve your validator pubkey.


### Create GenTx

Now that you have you key imported, you are able to use it to create your gentx.

To create the genesis transaction, you will have to choose the following parameters for your validator:

- moniker
- commission-rate
- commission-max-rate
- commission-max-change-rate
- min-self-delegation (must be >1)
- website (optional)
- details (optional)
- identity (keybase key hash, this is used to get validator logos in block explorers. optional)
- pubkey (gotten in previous step)

If you would like to override the memo field, use the `--ip` and `--node-id` flags.

An example genesis command would thus look like:

```sh
export node_name=<your node name>

celestia-appd keys add $node_name --keyring-backend=test

export node_addr=$(celestia-appd keys show $node_name -a --keyring-backend test)

echo $node_addr

celestia-appd add-genesis-account $node_addr "100000000000celestia" --keyring-backend test

celestia-appd gentx $node_name 5000000000celestia --keyring-backend=test --chain-id devnet-1
```

It will show an output something similar to:

```sh
Genesis transaction written to "/Users/ubuntu/.celestia-app/config/gentx/gentx-eb3b1768d00e66ef83acb1eee59e1d3a35cf76fc.json"
```

### Submit Your GenTx

To upload the your genesis file, please follow these steps:

1. Rename the gentx file just generated to gentx-{your-moniker}.json (please do not have any spaces or special characters in the file name)
2. Fork this repo by going to https://github.com/osmosis-labs/networks, clicking on fork, and choose your account (if multiple).
3. Clone your copy of the fork to your local machine
```sh
git clone https://github.com/<your_github_username>/networks
```
4. Copy the gentx to the networks repo (ensure that it is in the correct folder)

```sh
cp ~/.celestia-app/config/gentx/gentx-<your-moniker>.json networks/devnet-1/gentxs/
```

5. Commit and push to your repo.
 
```sh
cd networks
git add devnet-1/gentxs/*
git commit -m "<your validator moniker> gentx"
git push origin master
```

6. Create a pull request from your fork to master on this repo.

---
