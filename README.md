# Celestia Devnet
> Note: This guideline is only relevant for the current devnet. As we approach to testnet, there will be a new guide

- [Celestia Devnet](#celestia-devnet)
  - [Celestia 101 (Overview)](#celestia-101-overview)
    - [Celestia Full Node](#celestia-full-node)
    - [Celestia Light Node](#celestia-light-node)
  - [Installation guide](#installation-guide)
  - [Troubleshoot](#troubleshoot)

## Celestia 101 (Overview)
<i>Conceptually</i>, we have 2 main components in the network: 
1. `Celestia Full Node`.
2. `Celestia Light Node`.

> Note: When you see formatted text of those 2 above(`Celestia Full Node`, `Celestia Light Node`), remember that we are referring to a <i>conceptual</i> view. Standard formatted text is a technical view 

### Celestia Full Node
When we are mentioning `Celestia Full Node`, we are relating to a combination of 2 components: 

- Celestia Application (lives in [celestia-app](https://github.com/celestiaorg/celestia-app) repo)
- Celestia Full Node (lives in [celestia-node](https://github.com/celestiaorg/celestia-node) repo)

<i>Conceptually</i>, one can't live without another as we need every new block to be erasure coded in order to do data availability sampling from `Celestia Light Nodes`

Technically, Celestia Application and Celestia Full Node are operating on different processes of <b>one instance</b>(e.g. server). Let's break down what each of the parties do. 

Celestia Core in tandem with Celestia Application(C-App) are responsible for the following components:
- Consensus
- State Interaction
- Block Production
- Communication between C-App instances only!

Celestia Full Node(CFN) takes care of these points: 
- Sync blocks from C-App
- Erasure Codes the blocks
- Serve headers/shares to Celestia Light Nodes
- Communication between Celestia Full Nodes and Celestia Light Nodes

Communication of C-App and CFN processes happens over the RPC.

### Celestia Light Node
Celestia Light Nodes(CLN) lives only in [celestia-node](https://github.com/celestiaorg/celestia-node) repo. What does CLN do: 
1. Connects to a `Celestia Full Node`.
2. Syncs the headers.
3. Does data availability sampling (DAS) afterwards.

> Note: `Celestia Light Nodes` are not communicating between each other and only communicates between `Celestia Full Nodes`



## Installation guide
> Caveat 1: Make sure that you have at least 50+ Gb of free space to safely install+run the Celestia Application + Celestia Full Node and 5+ Gb of free space for Celestia Light Node

> Caveat 2: You need to install Celestia Application + Celestia Full Node in one HW instance and Celestia Light Node in the other HW instance

You need to install celestia-app, celestia-node in order to move on.
> Note: It is <b>much easier to run a non-validator Celestia App</b> + Celestia Full Node rather then going back and forth in order to be included in the active validator set

- [Celestia Application Installation](./xyz.md)
- [Celestia Node Installation](./celestia-node.md)

## Troubleshoot
Please navigate to [this link](./troubleshoot.md) to find more details


