# AGENTS.md

Configuration files (genesis, seeds, peers, addrbook) for Celestia networks, plus the scripts and workflows that maintain them. Node operators consume these files directly, so treat mainnet (`celestia/`) changes with care.

## Pull Requests

- Title format: `type(scope)!: [level] subject`, a conventional commit whose subject starts with exactly one of `[high]`, `[medium]`, `[low]` (e.g. `chore(peers): [medium] refresh mainnet peers`). See [README.md#contributing](README.md#contributing) for the level definitions.
- `[high]`: genesis files, chain IDs, new networks, hand-edited mainnet seeds or peers, security, or anything else a node operator consumes to join mainnet.
- `[medium]`: testnet config, automated peer refreshes, and script or workflow changes with limited impact.
- `[low]`: docs, tests, or tooling changes with low risk.
- Use the higher level when unsure. Don't repeat the review tag in the PR description; the title already carries it.

## Security PRs

- For PRs that resolve Hacken bug bounty reports, do not include details about the bug in the PR description. Instead, link to a Linear issue that contains more details on the bug and the link to the Hacken bug bounty report.
