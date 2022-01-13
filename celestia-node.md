# Running a Celestia Node

- [Running a Celestia Node](#running-a-celestia-node)
  - [Installation](#installation)
  - [Full Node Configuration](#full-node-configuration)
    - [Getting trusted hash](#getting-trusted-hash)
    - [Running full node](#running-full-node)
      - [1. Initialize the full node](#1-initialize-the-full-node)
      - [2. Edit configurations (adding other celestia full nodes)](#2-edit-configurations-adding-other-celestia-full-nodes)
  - [Light Client Configuration](#light-client-configuration)
    - [Getting trusted hash](#getting-trusted-hash-1)
    - [Running Light Client](#running-light-client)
- [Data Availability Sampling(DAS)](#data-availability-samplingdas)
  - [Pre-Requisites](#pre-requisites)
    - [Legend](#legend)
  - [Steps](#steps)

## Installation
Make sure that you have `git` and `golang` installed
```sh
git clone https://github.com/celestiaorg/celestia-node.git
make install
```


## Full Node Configuration

### Getting trusted hash
> Caveat: You need a running celestia-app in order to continue this guideline. Please refer to [celestia-app.md](https://github.com/celestiaorg/networks/celestia-app.md) for installation.


You need to have the trusted hash in order to initialize the Celestia full node
In order to know the hash, you need to query the celestia-app:
```sh
curl -s http://localhost:26657/block?height=1 | grep -A1 block_id | grep hash
```

### Running full node
#### 1. Initialize the full node
```sh
celestia full init --core.remote <ip:port of celestia-app> --headers.trusted-hash <hash_from_celestia_app>
```

Example:
```sh 
celestia full init --core.remote tcp://127.0.0.1:26657 --headers.trusted-hash 4632277C441CA6155C4374AC56048CF4CFE3CBB2476E07A548644435980D5E17
```

#### 2. Edit configurations (adding other celestia full nodes)

In order for your Celestia full node to communicate with other Celestia full nodes, then you need to add them as `mutual peers` in the `config.toml` file and allow the peer exchange.
```sh
nano ~/.celestia-full/config.toml
```
```sh
...
[P2P]
  ...
  #add multiaddresses of other celestia full nodes
  
  MutualPeers = [
    "/ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22", 
    "/ip4/x.x.x.x/tcp/yyy/p2p/abc"] #the /ip4/x.x.x.x is only for example. Don't add it! 
  PeerExchange = true #change this line to true. Be default it's false
  ...
...
```

1. Start the full node
```sh
celestia full start
```
Now, the Celestia full node will start syncing headers and storing blocks from Celestia application. 

> Note: At startup, we can see the `multiaddress` from Celestia full node. This is <b>needed for future light client</b> connections and communication between Celestia full nodes

Example:
```sh
/ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22
```

## Light Client Configuration

> Caveat: You don't need to run the light client on the same machine where Celestia full node is running

### Getting trusted hash
You need to have the trusted hash in order to initialize the light client
In order to know the hash, you need to query the celestia-app:

> Note: It is highly encouraged to run your own non-validating `celestia-app` node to get this trusted hash. However, you can ask for or take this hash from the discord/explorer if you want to just explore how easy it is to run the Celestia light client
```sh
curl -s http://<ip_address>:26657/block?height=1 | grep -A1 block_id | grep hash
``` 

### Running Light Client
> Note: If you want to run the light client only, then you can ask someone from the discord to send you the `multiaddress` from the Celestia full node to connect to

To start the light client, we need to know 2 variables:
- Trusted peer’s multi address to connect to (a Celestia full node is the case here)
- Trusted block hash from celestia-app

1. Initialize the light client

```sh
celestia full init --headers.trusted-peer <full_node_multiaddress> --headers.trusted-hash <hash_from_celestia_app>
```

Example: 

```sh 
celestia light init --headers.trusted-peer /ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22 --headers.trusted-hash 97682277DE3BA40176315102934EDB51CD9727FE31253C326F1F9581E14E2479
```

2. Start the light client
```sh
celestia light start
```
Now, the Celestia light client will start syncing headers. After sync is finished, light client will do data availability sampling(DAS) from the full node.

# Data Availability Sampling(DAS)

## Pre-Requisites
This is a list of runnining components you need in order to successfully continue this chapter:
- Celestia app
- Celestia light client

> Note: The light client should be connected to a Celestia full node to operate correctly. Either deploy your own Celestia full node or connect your <b>light client</b> to an existing Celestia full node in the network

### Legend
You will need 2 terminals in order to see how DASing works:
- First terminal(FT) run light client with logs info
- Second terminal(ST) submit payForMessage tx using celestia-app

## Steps
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
