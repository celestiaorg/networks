# Run a Celestia Light Node

- [Run a Celestia Light Node](#run-a-celestia-light-node)
  - [Dependencies](#dependencies)
    - [Update packages](#update-packages)
    - [Install GO](#install-go)
  - [Install Celestia Node](#install-celestia-node)
    - [Run the Light Node](#run-the-light-node)
  - [Data Availability Sampling(DAS)](#data-availability-samplingdas)
    - [Pre-Requisites](#pre-requisites)
    - [Create a wallet](#create-a-wallet)
    - [Fund the Wallet](#fund-the-wallet)
    - [Send a transaction](#send-a-transaction)
    - [Observe DAS in action](#observe-das-in-action)

## Dependencies
### Update packages
First, make sure to update and upgrade the OS:
```sh
sudo apt update && sudo apt upgrade -y
```
These are essential packages which are necessary execute many tasks like downloading files, compiling and monitoring the node:
```sh
sudo apt install curl tar wget clang pkg-config libssl-dev jq build-essential git make ncdu -y
```
### Install GO
It is necessary to install the GO language in the OS so we can later compile the Celestia Node. On our example, we are using version 1.17.2:
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

## Install Celestia Node
> Caveat 1: Make sure that you have at least 5+ Gb of free space for Celestia Light Node

Install the Celestia Node binary. Make sure that you have `git` and `golang` installed.
```sh
git clone https://github.com/celestiaorg/celestia-node.git
cd celestia-node/
make install
```

### Run the Light Node

> If you want to connect to your Celestia Bridge Node and start syncing the Celestia Light Node from a non-genesis hash, then consider editing `config.toml` file. 
More information on `config.toml`(./config-toml.md)

1. Initialize the Light Node
```sh
celestia light init
```

2. Start the Light Node

Start the Light Node as daemon process in the background
```sh
sudo tee <<EOF >/dev/null /etc/systemd/system/celestia-lightd.service
[Unit]
Description=celestia-lightd Light Node
After=network-online.target

[Service]
User=$USER
ExecStart=$HOME/go/bin/celestia light start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF
```

If the file was created succesfully you will be able to see its content:

```cat /etc/systemd/system/celestia-lightd.service```

3. Enable and start celestia-lightd daemon:

```sh
sudo systemctl enable celestia-lightd
sudo systemctl start celestia-lightd
```

4. Check if daemon has been started correctly:
```sh
sudo systemctl status celestia-lightd
```

5. Check daemon logs in real time:
```sh
sudo journalctl -u celestia-lightd.service -f
```

Now, the Celestia Light Node will start syncing headers. After sync is finished, Light Node will do data availability sampling(DAS) from the Bridge Node.

## Data Availability Sampling(DAS)

### Pre-Requisites
To continue, you will need:
- A Celestia Light Node connected to a Bridge Node
- A Celestia wallet

Open 2 terminals in order to see how DASing works:
1. First terminal: tail your Light Node logs
2. Second terminal: use Celestia App's CLI to submit a paid `payForMessage` tx to the network

### Create a wallet
First, you need a wallet to pay for the transaction.

**Option 1**: Use the Keplr wallet which has beta support for Celestia. https://staking.celestia.observer/

**Option 2**: Download the Celestia App binary which has a CLI for creating wallets
1. Download the celestia-appd binary inside `$HOME/go/bin` folder which will be used to create wallets.
```sh
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app/
make install
```
2. To check if the binary was succesfully compiled you can run the binary using the `--help` flag:
```sh
cd $HOME/go/bin
./celestia-appd --help
```

3. Create the wallet with any wallet name you want e.g.
```sh
celestia-appd keys add mywallet
```
Save the mnemonic output as this is the only way to recover your validator wallet in case you lose it! 

### Fund the Wallet
You can fund an existing wallet via Discord by sending this message to #faucet channel:
```
!faucet celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Wait to see if you get a confirmation that the tokens have been successfully sent. To check if tokens have arrived succesfully to the destination wallet run the command below replacing the public address with your own:
```sh
celestia-appd q bank balances celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Send a transaction
In the second terminal, submit a `payForMessage` transaction with `celestia-appd` (or do so in the wallet):
```sh
celestia-appd tx payment payForMessage <hex_namespace> <hex_message> --from <wallet_name> --keyring-backend <keyring-name> --chain-id <chain_name>
```
Example:
```sh 
celestia-appd tx payment payForMessage 0102030405060708 68656c6c6f43656c6573746961444153 --from myWallet --keyring-backend test --chain-id devnet-2
```

### Observe DAS in action
In the Light Node logs you should see how data availability sampling works:

Example:
```sh
INFO	das	das/daser.go:96	sampling successful	{"height": 81547, "hash": "DE0B0EB63193FC34225BD55CCD3841C701BE841F29523C428CE3685F72246D94", "square width": 2, "finished (s)": 0.000117466}
```
