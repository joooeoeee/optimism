// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Script } from "forge-std/Script.sol";
import { Artifacts } from "scripts/Artifacts.s.sol";
import { Deploy } from "scripts/deploy/Deploy.s.sol";
import { Config } from "scripts/Config.sol";
import { L2Genesis } from "scripts/L2Genesis.s.sol";
import { console } from "forge-std/console.sol";

/// @title InteropDeployer
/// @notice Functions to deploy OP-Stack components in incremental steps, for interop testing usage.
contract InteropDeploy is Script {

    L2Genesis l2GenesisUtils = L2Genesis(address(uint160(uint256(keccak256(abi.encode("optimism.l2genesis"))))));

    Deploy deployUtils = Deploy(address(uint160(uint256(keccak256(abi.encode("optimism.deploy"))))));

    /// @notice The address of the deployer account.
    address internal deployer;

    uint80 internal constant DEV_ACCOUNT_FUND_AMT = 10_000_000 ether;

    function setUp() public virtual override {
        deployer = makeAddr("deployer");

        // Note: we don't init this contract, we just use some helper functions
        vm.etch(address(l2GenesisUtils), vm.getDeployedCode("L2Genesis.s.sol:L2Genesis"));
        vm.label(address(l2GenesisUtils), "L2Genesis");
        vm.allowCheatcodes(address(l2GenesisUtils));

        vm.etch(address(deployUtils), vm.getDeployedCode("Deploy.s.sol:Deploy"));
        vm.label(address(deployUtils), "Deploy");
        vm.allowCheatcodes(address(deployUtils));
    }

    function loadAllocs() public {
        string memory inputAllocsPath = vm.envString("ALLOCS_INPUT_PATH");
        vm.loadAllocs(inputAllocsPath);
    }

    function dumpAllocs() public {
        /// Reset so its not included state dump
        vm.etch(address(l2GenesisUtils), "");

        vm.etch(msg.sender, "");
        vm.resetNonce(msg.sender);
        vm.deal(msg.sender, 0);

        vm.deal(deployer, 0);
        vm.resetNonce(deployer);

        string memory outputAllocsPath = vm.envString("ALLOCS_OUTPUT_PATH");
        vm.dumpState(outputAllocsPath);
    }

    function initialL1() public {
        // TODO set chain ID to L1 from deploy config (set via script for now?)

        // Put all preinstalls into the L1 chain.
        // This includes the 4788 beacon roots contract.
        // Hack: reuse the L2 genesis flow;
        // to be refactored into common util later.
        l2GenesisUtils.setPreinstalls();

        // prefund user dev accounts
        l2GenesisUtils.fundDevAccounts();

        // TODO: loading these addresses from a well encapsulated dev config would be nice.
        // operator accounts
        vm.deal(address(0xCA55aC8514b25C660151a8AE0c90f116DF160daa), DEV_ACCOUNT_FUND_AMT); // batcher A
        vm.deal(address(0x252a3336Fb2A4352D1bD3b139f4e540AA45236bd), DEV_ACCOUNT_FUND_AMT); // batcher B
        vm.deal(address(0x71b4a2d9B91726bdb5849D928967A1654D7F3de7), DEV_ACCOUNT_FUND_AMT); // proposer A
        vm.deal(address(0x8c408c9ce6718F4a3AFa7860f2E7B190B25fBDfA), DEV_ACCOUNT_FUND_AMT); // proposer B
        // TODO proxy admins, sys config owners

        dumpAllocs();
    }

    function deploySuperchain() public {
        loadAllocs();

        // TODO load config

        // TODO set chain ID to L1 from deploy config
        // TODO load superchain deploy config

        // Deploy superchain components;
        // AddressManager, ProxyAdmin, SuperchainConfig, ProtocolVersions
        deployUtils.setupSuperchain();

        // Deploy the implementations; to be reused between L2s (OP-Stack Manager style)
        deployUtils.deployImplementations();

        dumpAllocs();
    }

    function deployL2() public {
        loadAllocs();

        // TODO set chain ID to L1 from deploy config
        // TODO load L2 deploy config

        deployUtils.deploySafe("SystemOwnerSafe");
        // TODO need to prepare artifacts such that the proxies
        // can hook to the superchain implementations
        deployUtils.deployProxies();
        deployUtils.initializeImplementations();

        // For some reason the FP system deployment is separate
        // from the above MCP-like proxy/impl setup.
        deployUtils.setAlphabetFaultGameImplementation({ _allowUpgrade: false });
        deployUtils.setFastFaultGameImplementation({ _allowUpgrade: false });
        deployUtils.setCannonFaultGameImplementation({ _allowUpgrade: false });
        deployUtils.setPermissionedCannonFaultGameImplementation({ _allowUpgrade: false });
        deployUtils.transferDisputeGameFactoryOwnership();
        deployUtils.transferDelayedWETHOwnership();

        dumpAllocs();
    }
}
