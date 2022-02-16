# Celestia Devnet
- [Celestia Devnet](#celestia-devnet)
  - [Celestia 101 (Overview)](#celestia-101-overview)
    - [Celestia Full Node](#celestia-full-node)
    - [Celestia Light Node](#celestia-light-node)
  - [Installation guide](#installation-guide)
  - [Troubleshoot](#troubleshoot)

> Note: This guideline is only relevant for the current devnet. As we approach to testnet, there will be a new guide

## Quickstart

Devnet participants have the option of running:

1. [Celestia Light Nodes](#celestia-light-nodes) (low CPU, 5GB+ disk): [get started](celestia-light-node.md)
2. [Bridge Nodes](#bridge-nodes) (low CPU, 100GB+ disk): [get started](celestia-bridge-node.md)
3. [Bridge Validator Nodes](#bridge-validator-nodes) (same node as 2): [get started](celestia-bridge-node.md#running-a-validating-bridge-node)
4. _Celestia Full DA Nodes are under development_: [see ADR](https://github.com/celestiaorg/celestia-node/blob/main/docs/adr/adr-003-march2022-testnet.md#full-node)

You can view chain activity at the current `devnet-2` explorer: https://celestia.observer/
## Overview

Devnet demonstrates Celestia’s data availability capabilities on a libp2p network and live blockchain. 

The current Devnet implementation runs two individual but connected networks:

1. A libp2p network with [**Light Nodes**](#celestia-light-nodes), that handles data availability interactions
2. A p2p network with [**Validators Nodes**](#bridge-validator-nodes), that handles the underlying consensus and block production. 

Special [**Bridge Nodes**](#bridge-nodes) process blocks from the underlying consensus network to the data availability network.

![Network Overview](diagrams/NetworkOverview.png)

> It’s important to note that mainnet may look very different from this devnet implementation, as the architecture continues to be improved. You can read more about devnet decisions here and here (link to ADR). 

---

## Celestia Light Nodes

Light nodes (CLN) ensure data availability and can publish transactions. This is the most ubiquitous way to interact with the Celestia network.

Specifically, Light Nodes: 

1. Connect to a `Celestia Bridge Node` in the DA network. *Note: Light Nodes do not communicate with each other, but only with Bridge Nodes.*
2. Listen for `ExtendedHeader`s, i.e. wrapped “raw” headers, that notify Celestia Nodes of new block headers and relevant DA metadata.
3. Perform data availability sampling (DAS) on the received headers

![Light Nodes](diagrams/LightNodes.png)

**Installation**
- [Light Node Setup Guide](/celestia-light-node.md)
- Source code repository(s):
    - [celestia-node](https://github.com/celestiaorg/celestia-node)

---

## Bridge Nodes

Bridge Nodes connect the aforementioned p2p and libp2p networks.

Specifically, Bridge Nodes: 

1. Import and process “raw” headers & blocks from a trusted Core process in the p2p network. *Note: Celestia Core can be run as either an internal process or simply accessed via a remote endpoint.* 
2. Validate and erasure code the blocks
3. Produce block shares to Light Nodes in the DA network.

![Bridge Nodes](diagrams/BridgeNodes.png)

From an implementation perspective, Bridge Nodes run two separate processes: 

1. Celestia App (with a Celestia Core) [see repo](https://github.com/celestiaorg/celestia-app)
    - **Celestia app** is the state machine where the application and the proof-of-stake logic is run. 
    > Celestia App is built on [Cosmos SDK](https://docs.cosmos.network/) and also encompasses **Celestia Core**.
    - **Celestia Core** is the state interaction, consensus and block production layer. 
    > Celestia Core is built on [Tendermint Core](https://docs.tendermint.com/), modified to store (1) invalid transactions and (2) data roots of erasure coded blocks, among other changes. see ADR.
2. Celestia Node [see repo](https://github.com/celestiaorg/celestia-node)
    - **Celestia Node** augments the above with a separate libp2p network that serves data availability sampling requests. The team sometimes refer to this as the "halo" network.

**Installation**
- [Bridge Node Setup Guide](/celestia-bridge-node.md)
- Repositories:
    - [celestia-app](https://github.com/celestiaorg/celestia-app)
    - [celestia-node](https://github.com/celestiaorg/celestia-node)

---

## Bridge Validator Nodes

Bridge Nodes have the option of validating the P2P network using its Celestia App component. However, running validator nodes is not a requirement to learn about Celestia’s main value proposition.

> Only the top 100 validators make it into the active validator set. The team is not looking for additional validators at the moment and recommend running light or full nodes instead.

**Installation**

- [Validator Setup Guide](celestia-bridge-node.md#running-a-validating-bridge-node)
- Repositories:
    - [celestia-app](https://github.com/celestiaorg/celestia-app)

## Troubleshoot
Please navigate to [this link](./troubleshoot.md) to find more details


