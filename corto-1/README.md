# Corto

Corto is an internal testnet with chain ID `corto-1`. This genesis starts
directly at app version 10. It preserves the previous genesis allocations
and signed gentx; only `consensus.params.version.app` changes from 9 to 10.

## Genesis checksum

SHA-256 of `genesis.json`:

```text
4325c418ae9bd02e5b1c87aef2805243096a1951e7908b464675623f13cba1d2
```

This is the file checksum, not the block-1 hash used by DA nodes. Capture the
new block-1 hash after the coordinated restart.

## Coordinated reset

Existing nodes must reset their chain state before using this genesis.
Changing the genesis file over existing databases is not an upgrade. Stop
all validators and signers before resetting state, preserve their identities,
and distribute identical genesis bytes to every consensus node. DA nodes
and downstream clients must also discard old-history stores and trust anchors.

The genesis contains only validator-0's gentx. Validators registered after
genesis must rejoin on-chain; their existing keys can be reused.

The `v10.2.0-corto` runtime supports this genesis, but its `download-genesis`
command pins the previous file checksum and will reject the new file. Until
a release with the updated checksum is installed, distribute this file
directly and verify its SHA-256.

Do not publish this replacement until the coordinated reset window.
