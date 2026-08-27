# Celestia Networks

This repository contains the configuration files for Celestia networks.

Please refer to the [Celestia Docs](https://docs.celestia.org) for guides on running your own node. The configuration files `genesis.json` and `addrbook.json` are intended for use by Consensus nodes and should be placed in the node's home directory (by default `$HOME/.celestia-app/config`).

## Seeds and peers

Some networks include `seeds.txt` and/or `peers.txt` files. Each line is a node address in the form `nodeID@host:port` (for example `8084e73b70dbe7fba3602be586de45a516012e6f@144.76.112.238:26656`).

A brand new node has no way to find the rest of the network on its own, so it needs at least one known address to bootstrap from. These files provide those starting points:

- **Seeds** (`seeds.txt`) are nodes that exist purely for peer discovery. Your node connects to a seed, downloads a list of other addresses in the network, and then disconnects. Seeds do not maintain a persistent connection, so they are the lightweight way to discover peers when joining a network.
- **Peers** (`peers.txt`) are nodes you connect to and stay connected to (configured as persistent peers). Unlike seeds, your node keeps these connections open to relay blocks and transactions.

Node operators need these because a node cannot sync the chain or participate in consensus until it has connected to the network. Provide them via your node's config or command-line flags:

- Add the contents of `seeds.txt` to the `seeds` field in `config.toml` (or pass `--p2p.seeds`).
- Add the contents of `peers.txt` to the `persistent_peers` field in `config.toml` (or pass `--p2p.persistent_peers`).

Multiple entries are comma-separated when placed in config or flags. See the [Celestia Docs](https://docs.celestia.org) for the exact steps for the node type you are running.

## Networks

| Name     | Type         | Chain ID     | Configs                    |
|----------|--------------|--------------|----------------------------|
| Mocha    | Testnet      | `mocha-4`    | [mocha-4](./mocha-4)       |
| Arabica  | Testnet      | `arabica-11` | [arabica-11](./arabica-11) |
| Corto    | Testnet      | `corto-1`    | [corto-1](./corto-1)       |
| Celestia | Mainnet Beta | `celestia`   | [celestia](./celestia)     |

## Software versions

The celestia-app and celestia-node versions in use on each network are visible [in the Celestia Docs](https://docs.celestia.org/operate/networks/overview/).
