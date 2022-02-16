# Running a Celestia Light Node

- [Running a Celestia Node](#running-a-celestia-node)
  - [Installation](#installation)
  - [Light Node Configuration](#light-node-configuration)
    - [Getting trusted hash](#getting-trusted-hash-1)
    - [Running Light Node](#running-light-node)
- [Data Availability Sampling (DAS)](#data-availability-samplingdas)
  - [Pre-Requisites](#pre-requisites)
    - [Legend](#legend)
  - [Steps](#steps)

## Installation
Make sure that you have `git` and `golang` installed
```sh
git clone https://github.com/celestiaorg/celestia-node.git
make install
```

## Light Node Configuration

Light Nodes must reference a trusted hash and connect to a trusted Bridge Node `multiaddress` to initialize. 

If you just want to explore how easy it is to run a Celestia Light Node, you can use the constants provided by the following guide. 

For added security, you can run your own Bridge Node (guide here)[TODO]. Note that you don't need to run the Light Node on the same machine.

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
#### Option 1: Revisit daemon logs for the its IP4 address
```
TODO: get the first 50 logs from Bridge Node daemon
```
#### Option 2: Use [this multiaddress](/devnet-2/celestia-node/mutual_peers.txt)

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
Now, the Celestia Light Node will start syncing headers. After sync is finished, Light Node will do data availability sampling(DAS) from the full node.

# Data Availability Sampling(DAS)

## Pre-Requisites
This is a list of runnining components you need in order to successfully continue this chapter:
- Celestia Light Node
- Light Node Connection to a Bridge Node

> Note: The Light Node should be connected to a Bridge Node to operate correctly. Either deploy your own Bridge Node or connect your <b>Light Node</b> to an existing Bridge Node in the network

### Legend
You will need 2 terminals in order to see how DASing works:
- First terminal(FT) run Light Node with logs info
- Second terminal(ST) submit payForMessage tx using celestia-app

## Steps
1. In (ST) Submit a `payForMessage` transaction with `celestia-appd`
```sh
celestia-appd tx payment payForMessage <hex_namespace> <hex_message> --from <node_name> --keyring-backend <keyring-name> --chain-id <chain_name>
```
Example:
```sh 
celestia-appd tx payment payForMessage 0102030405060708 68656c6c6f43656c6573746961444153 --from eva00 --keyring-backend test --chain-id devnet-2
```
2. In (FT) you should see in logs how DAS is working

Example:
```sh
INFO	das	das/daser.go:96	sampling successful	{"height": 81547, "hash": "DE0B0EB63193FC34225BD55CCD3841C701BE841F29523C428CE3685F72246D94", "square width": 2, "finished (s)": 0.000117466}
```
