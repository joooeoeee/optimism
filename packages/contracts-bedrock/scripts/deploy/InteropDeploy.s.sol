// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Script } from "forge-std/Script.sol";
import { Artifacts } from "scripts/Artifacts.s.sol";
import { Deploy } from "scripts/deploy/Deploy.s.sol";
import { Config } from "scripts/Config.sol";
import { console } from "forge-std/console.sol";

/// @title InteropDeployer
/// @notice Functions to deploy OP-Stack components in incremental steps, for interop testing usage.
contract InteropDeploy is Deploy {

    function initialL1() public {
        // TODO load L1 dev deploy config
        // TODO set chain ID to L1 from deploy config
        // TODO add all preinstalls
        // TODO add 4788 beacon roots contract
        // TODO prefund dev accounts
        // TODO prefund batcher/proposer accounts of chain A and B
        // TODO dump allocs
    }

    function deploySuperchain() public {
        // TODO load allocs
        // TODO set chain ID to L1 from deploy config
        // TODO load superchain deploy config
        // TODO deploy proxy admin
        // TODO deploy SuperchainConfig
        // TODO deploy ProtocolVersions
        // TODO dump allocs
    }

    function deployL2() public {
        // TODO load allocs
        // TODO set chain ID to L1 from deploy config
        // TODO load L2 deploy config
        // TODO deploy proxies, attached to superchain implementation contracts
        // TODO dump allocs
    }
}
