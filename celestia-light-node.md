# Running a Celestia Light Node

- [Running a Celestia Node](#running-a-celestia-node)
  - [Dependencies](#dependencies)
  - [Install Celestia Node](#install-celestia-node)
  - [Light Node Configuration](#light-node-configuration)
    - [Getting trusted hash](#getting-trusted-hash-1)
    - [Running Light Node](#running-light-node)
- [Data Availability Sampling (DAS)](#data-availability-samplingdas)
  - [Pre-Requisites](#pre-requisites)
    - [Legend](#legend)
  - [Steps](#steps)

> Caveat 1: Make sure that you have at least 5+ Gb of free space for Celestia Light Node
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
### Installing GO
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

## Install Celestia Node
Make sure that you have `git` and `golang` installed
```sh
git clone https://github.com/celestiaorg/celestia-node.git
make install
```

## Light Node Configuration

Light Nodes must reference a trusted hash and connect to a trusted Bridge Node `multiaddress` to initialize. 

If you just want to explore how easy it is to run a Celestia Light Node, you can use the constants provided by the following guide. 

> For added security, you can run your own Bridge Node ([guide here](/celestia-bridge-node.md)). Note that you don't need to run the Light Node on the same machine.

### Getting the trusted hash
You need to have the trusted hash in order to initialize the Light Node

#### Option 1: Query your Bridge Node for the trusted hash:
```sh
curl -s http://<ip_address>:26657/block?height=1 | grep -A1 block_id | grep hash
``` 

#### Option 2: Use the following hash provided by the team: 
```sh
4632277C441CA6155C4374AC56048CF4CFE3CBB2476E07A548644435980D5E17
```

### Getting the trusted multiaddress
#### Option 1: Get IP4 from your Bridge Node's daemon logs
```
journalctl -u YOUR_CELESTIA_NODE.service --since "NODE_START_TIME" --until "1_MIN_AFTER_START_TIME"
```
#### Option 2: Use [this multiaddress](/devnet-2/celestia-node/mutual_peers.txt) provided in this repo

### Running Light Node
1. Initialize the Light Node

```sh
celestia full init --headers.trusted-peer <full_node_multiaddress> --headers.trusted-hash <hash_from_celestia_app>
```

Example: 

```sh 
celestia light init --headers.trusted-peer /ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22 --headers.trusted-hash 4632277C441CA6155C4374AC56048CF4CFE3CBB2476E07A548644435980D5E17
```

2. Start the Light Node
```sh
celestia light start
```

Or start it as a daemon process in the background
```sh
sudo tee <<EOF >/dev/null /etc/systemd/system/celestia-lightd.service
[Unit]
Description=celestia-lightd LightNode daemon
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

Enable and start celestia-lightd daemon:

```sh
sudo systemctl enable celestia-lightd
sudo systemctl start celestia-lightd
```

Check if daemon has been started correctly:
```sh
sudo systemctl status celestia-lightd
```

Check daemon logs in real time:
```sh
sudo journalctl -u celestia-lightd.service -f
```

Now, the Celestia Light Node will start syncing headers. After sync is finished, Light Node will do data availability sampling(DAS) from the full node.

# Data Availability Sampling(DAS)

## Pre-Requisites
This is a list of runnining components you need in order to successfully continue this chapter:
- Celestia Light Node
- Light Node Connection to a Bridge Node
- A Celestia wallet 

> Note: The Light Node should be connected to a Bridge Node to operate correctly. Either deploy your own Bridge Node or connect your <b>Light Node</b> to an existing Bridge Node in the network

### Legend
You will need 2 terminals in order to see how DASing works:
- First terminal(FT) run Light Node with logs info
- Second terminal(ST) submit payForMessage tx using celestia-app

## Steps
1. Create a Celestia wallet. 
Download the celestia-appd binary inside `$HOME/go/bin` folder which will be used to create wallets.
```sh
git clone https://github.com/celestiaorg/celestia-app.git
cd celestia-app/
make install
```
To check if the binary was succesfully compiled you can run the binary using the `--help` flag:
```sh
cd $HOME/go/bin
./celestia-appd --help
```

Create the wallet with any wallet name you want e.g.
```sh
celestia-appd keys add mywallet
```
Save the mnemonic output as this is the only way to recover your validator wallet in case you lose it! 

You can fund an existing wallet via Discord by sending this message to #faucet channel:
```
!faucet celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
Wait to see if you get a confirmation that the tokens have been successfully sent. To check if tokens have arrived succesfully to the destination wallet run the command below replacing the public address with your own:
```sh
celestia-appd q bank balances celes1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

2. In (ST) Submit a `payForMessage` transaction with `celestia-appd`
```sh
celestia-appd tx payment payForMessage <hex_namespace> <hex_message> --from <wallet_name> --keyring-backend <keyring-name> --chain-id <chain_name>
```
Example:
```sh 
celestia-appd tx payment payForMessage 0102030405060708 68656c6c6f43656c6573746961444153 --from myWallet --keyring-backend test --chain-id devnet-2
```

3. In (FT) you should see in logs how DAS is working

Example:
```sh
INFO	das	das/daser.go:96	sampling successful	{"height": 81547, "hash": "DE0B0EB63193FC34225BD55CCD3841C701BE841F29523C428CE3685F72246D94", "square width": 2, "finished (s)": 0.000117466}
```