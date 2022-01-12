# Run Celestia Devnet
> Note: This guideline is only relevant for the current devnet. As we approach to testnet, there will be a new guide

- [Run Celestia Devnet](#run-celestia-devnet)
  - [Celestia 101(Overview)](#celestia-101overview)
    - [Celestia Full Node](#celestia-full-node)
    - [Celestia Light Client](#celestia-light-client)
  - [Installation guide](#installation-guide)
  - [Troubleshoot](#troubleshoot)

## Celestia 101(Overview)
<i>Conceptually</i>, we have 2 main components in the network: 
1. `Celestia Full Node`.
2. `Celestia Light Client`.

> Note: When you see formated text of those 2 above(`Celestia Full Node`, `Celestia Light Client`), remember that we are refering to a <i>conceptual</i> view. Standard formated text is a technical view 

### Celestia Full Node
When we are mentioning `Celestia Full Node`, we are relating to a combination of 2 components: 

- Celestia Application (lives in [celestia-app](https://github.com/celestiaorg/celestia-app) repo)
- Celestia Full Node (lives in [celestia-node](https://github.com/celestiaorg/celestia-node) repo)

<i>Conceptually</i>, one can't live without another as we need every new block to be erasure coded in order to do data availability sampling from `Celestia Light Clients`

Technically, Celestia Application and Celestia Full Node are operating on different processes of <b>one instance</b>(e.g. server). Let's break down what each of the parties do. 

Celestia Application(CA) takes care of these points: 
- Consensus
- State Interaction
- Block Production
- Communication between Celestia Application instances only!

Celestia Full Node(CFN) takes care of these points: 
- Sync blocks from Celestia App
- Erasure Codes the blocks
- Serve headers/shares to Celestia Light Clients
- Communication between Celestia Full Nodes and Celestia Light Clients

Communication of CA and CFN processes happens over the RPC.

### Celestia Light Client
Celestia Light Client(CLC) lives only in [celestia-node](https://github.com/celestiaorg/celestia-node) repo. What does CLC do: 
1. Connects to a `Celestia Full Node`.
2. Syncs the headers.
3. Does data availability sampling(DAS) afterwards.

> Note: `Celestia Light Clients` are not communicating between each other and only communicates between `Celestia Full Nodes`



## Installation guide
You need to install celestia-app, celestia-node in order to move on.
> Note: It is <b>much easier to run a non-validator celestia app</b> + celestia full node rather then going back and forth in order to be included in the active validator set

- [Celestia Application Installation](./xyz.md)
- [Celestia Node Installation](./celestia-node.md)

## Troubleshoot
Please navigate to [this link](./troubleshoot.md) to find more details


