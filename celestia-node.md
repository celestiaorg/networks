# Running a Celestia Node

- [Running a Celestia Node](#running-a-celestia-node)
  - [Installation](#installation)
  - [Understanding config.toml](#understanding-configtoml)
    - [[Core]](#core)
    - [[P2P]](#p2p)
      - [Bootstrap](#bootstrap)
      - [Mutual Peers](#mutual-peers)
    - [[Services]](#services)
      - [TrustedHash and TrustedPeer](#trustedhash-and-trustedpeer)
  - [Full Node Configuration](#full-node-configuration)
    - [Getting trusted hash](#getting-trusted-hash)
    - [Running full node](#running-full-node)
      - [1. Initialize the full node](#1-initialize-the-full-node)
      - [2. Edit configurations (adding other celestia full nodes)](#2-edit-configurations-adding-other-celestia-full-nodes)
  - [Light Node Configuration](#light-node-configuration)
    - [Getting trusted hash](#getting-trusted-hash-1)
    - [Running Light Node](#running-light-node)
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

## Understanding config.toml
After initialization, for any type of node, you will find a `config.toml`. Let's break down some of the most used sections
### [Core]
This section is needed for the Celestia Full Node. By default, `Remote = false`. Still for devnet, we are going to use the remote core option and this can also be set
by the command line flag `--core.remote` 

### [P2P]
#### Bootstrap
By default, the `Bootstrapper = false` and the `BootstrapPeers` is empty. You do need to activate this section and add bootstrap peers in order to bootstrap your node

#### Mutual Peers
The purpose of this config is to set up a bidirectional communication. This is usually the case for Celestia Full Nodes. In addition, you need to change the field 
`PeerExchange` from false to true

### [Services]
#### TrustedHash and TrustedPeer
TrustedHash is needed to properly initialize Celestia Full Node with the active `Remote` Core Client. Celestia Light Node needs to be initialized with the trusted hash, too.
TrustedPeer is the most crucial field to be filled for correct Celestia Light Node execution. Any Celestia Full Node can be a trusted peer for the Light one. However, the Light node
by design can not be a trusted peer for another Light Node.


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

In order for your Celestia full node to communicate with other Celestia full nodes, then you need to add them as `mutual peers` in the `config.toml` file and allow the peer exchange. Please navigate to `networks/devnet-2/celestia-node/mutual_peers.txt` to find the list of mutual peers
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

> Note: At startup, we can see the `multiaddress` from Celestia full node. This is <b>needed for future Light Node</b> connections and communication between Celestia full nodes

Example:
```sh
/ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22
```

## Light Node Configuration

> Caveat: You don't need to run the Light Node on the same machine where Celestia full node is running

### Getting trusted hash
You need to have the trusted hash in order to initialize the Light Node
In order to know the hash, you need to query the celestia-app:

> Note: It is highly encouraged to run your own non-validating `celestia-app` node to get this trusted hash. However, you can ask for or take this hash from the discord/explorer if you want to just explore how easy it is to run the Celestia Light Node
```sh
curl -s http://<ip_address>:26657/block?height=1 | grep -A1 block_id | grep hash
``` 

### Running Light Node
> Note: If you want to run the Light Node only, then you can ask someone from the discord to send you the `multiaddress` from the Celestia full node to connect to

To start the Light Node, we need to know 2 variables:
- Trusted peer’s multi address to connect to (a Celestia full node is the case here)
- Trusted block hash from celestia-app

1. Initialize the Light Node

```sh
celestia full init --headers.trusted-peer <full_node_multiaddress> --headers.trusted-hash <hash_from_celestia_app>
```

Example: 

```sh 
celestia light init --headers.trusted-peer /ip4/46.101.22.123/tcp/2121/p2p/12D3KooWD5wCBJXKQuDjhXFjTFMrZoysGVLtVht5hMoVbSLCbV22 --headers.trusted-hash 97682277DE3BA40176315102934EDB51CD9727FE31253C326F1F9581E14E2479
```

2. Start the Light Node
```sh
celestia light start
```
Now, the Celestia Light Node will start syncing headers. After sync is finished, Light Node will do data availability sampling(DAS) from the full node.

# Data Availability Sampling(DAS)

## Pre-Requisites
This is a list of runnining components you need in order to successfully continue this chapter:
- Celestia app
- Celestia Light Node

> Note: The Light Node should be connected to a Celestia full node to operate correctly. Either deploy your own Celestia full node or connect your <b>Light Node</b> to an existing Celestia full node in the network

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
