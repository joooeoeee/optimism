#!/bin/bash

set -eu

# Run this with workdir set as root of the repo
if [ -f "versions.json" ]; then
    echo "Running create-chains script."
else
    echo "Cannot run create-chains script, must be in root of repository, but currently in:"
    echo "$(pwd)"
    exit 1
fi

# Check if already created
if [ -d ".devnet-interop" ]; then
    echo "Already created chains."
    exit 1
else
    echo "Creating new interop devnet chain configs"
fi

mkdir ".devnet-interop"

# generate L1 allocs with dev-allocs (incl. accounts for chain A and B proposer/batcher!)
cd ./packages/contracts-bedrock && \
  DEPLOYMENT_OUTFILE="../../.devnet-interop/deployments-l1-dev.json" \
  DEPLOY_CONFIG_PATH="../../interop-devnet/deploy-config-l1-dev.json" \
  ALLOCS_OUTPUT_PATH="../../.devnet-interop/allocs-l1-dev.json" \
  forge script "./scripts/deploy/InteropDeploy.s.sol:InteropDeploy" \
    --sig "initialL1()" \
    --sender "0x90F79bf6EB2c4f870365E785982E1f101E93b906"

# deploy superchain to L1
cd ./packages/contracts-bedrock && \
  DEPLOYMENT_OUTFILE="../../.devnet-interop/deployments-l1-superchain.json" \
  DEPLOY_CONFIG_PATH="../../interop-devnet/deploy-config-superchain.json" \
  ALLOCS_INPUT_PATH="../../.devnet-interop/allocs-l1-dev.json" \
  ALLOCS_OUTPUT_PATH="../../.devnet-interop/allocs-l1-superchain.json" \
  forge script "./scripts/deploy/InteropDeploy.s.sol:InteropDeploy" \
    --sig "superchainL1()" \
    --sender "0x90F79bf6EB2c4f870365E785982E1f101E93b906"

# deploy L1-contracts A
cd ./packages/contracts-bedrock && \
  DEPLOYMENT_OUTFILE="../../.devnet-interop/deployments-l1-a.json" \
  DEPLOY_CONFIG_PATH="../../interop-devnet/deploy-config-l2-a.json" \
  ALLOCS_INPUT_PATH="../../.devnet-interop/allocs-l1-superchain.json" \
  ALLOCS_OUTPUT_PATH="../../.devnet-interop/allocs-l1-superchain-and-l2-a.json" \
  forge script "./scripts/deploy/InteropDeploy.s.sol:InteropDeploy" \
    --sig "deployL2()" \
    --sender "0x90F79bf6EB2c4f870365E785982E1f101E93b906"

# deploy L1-contracts B
cd ./packages/contracts-bedrock && \
  DEPLOYMENT_OUTFILE="../../.devnet-interop/deployments-l1-b.json" \
  DEPLOY_CONFIG_PATH="../../interop-devnet/deploy-config-l2-b.json" \
  ALLOCS_INPUT_PATH="../../.devnet-interop/allocs-l1-superchain-and-l2-a.json" \
  ALLOCS_OUTPUT_PATH="../../.devnet-interop/allocs-l1-complete.json" \
  forge script "./scripts/deploy/InteropDeploy.s.sol:InteropDeploy" \
    --sig "deployL2()" \
    --sender "0x90F79bf6EB2c4f870365E785982E1f101E93b906"

# create L2 A allocs
cd ./packages/contracts-bedrock && \
  CONTRACT_ADDRESSES_PATH="../../.devnet-interop/deployments-l1-a.json" \
  DEPLOY_CONFIG_PATH="../../interop-devnet/deploy-config-l2-a.json" \
  forge script "./scripts/L2Genesis.s.sol:L2Genesis" --sig "runWithAllUpgrades()"

# create L2 B allocs
cd ./packages/contracts-bedrock && \
  CONTRACT_ADDRESSES_PATH="../../.devnet-interop/deployments-l1-b.json" \
  DEPLOY_CONFIG_PATH="../../interop-devnet/deploy-config-l2-b.json" \
  forge script "./scripts/L2Genesis.s.sol:L2Genesis" --sig "runWithAllUpgrades()"

# create L1 EL genesis
# TODO this is all kinds of broken; fix l1-deployments dependency
#go run ./op-node/cmd genesis l1 \
#  --deploy-config paths.devnet_config_path \
#  --l1-allocs paths.allocs_l1_path \
#  --l1-deployments paths.addresses_json_path \
#  --outfile.l1 paths.genesis_l1_path

# create L2 CL genesis
eth2-testnet-genesis deneb \
  --config=./beacon-data/config.yaml \
  --preset-phase0=minimal \
  --preset-altair=minimal \
  --preset-bellatrix=minimal \
  --preset-capella=minimal \
  --preset-deneb=minimal \
  --eth1-config=../.devnet-interop/genesis-l1.json \
  --state-output=../.devnet-interop/genesis-l1.ssz \
  --tranches-dir=../.devnet-interop/tranches \
  --mnemonics=mnemonics.yaml \
  --eth1-withdrawal-address=0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
  --eth1-match-genesis-time

# create L2 A genesis + rollup config
# TODO this needs to be refactored
# Also the L1 RPC part should be static, to avoid temporary L1 node going up.
#go run ./op-node/cmd/main.go genesis l2 \
#  --l1-rpc http://localhost:8545 \
#  --deploy-config devnet_config_path \
#  --l2-allocs ./.devnet-interop/allocs-l2-a.json \
#  --l1-deployments ./.devnet-interop/l1-deployments-a.json \
#  --outfile.l2 ./.devnet-interop/genesis-l2-a.json \
#  --outfile.rollup ./.devnet-interop/rollup-a.json

# create L2 B genesis + rollup config
# TODO repeat for L2 B
