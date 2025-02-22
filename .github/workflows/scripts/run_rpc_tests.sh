#!/bin/bash

set +e # Disable exit on error

# Array of disabled tests
disabled_tests=(
    # Failing after the PR https://github.com/erigontech/erigon/pull/13903 - diff is only an error message in the result
    eth_estimateGas/test_14.json
    # Failing after the PR https://github.com/erigontech/erigon/pull/13617 that fixed this incompatibility
    # issues https://hive.pectra-devnet-5.ethpandaops.io/suite.html?suiteid=1738266984-51ae1a2f376e5de5e9ba68f034f80e32.json&suitename=rpc-compat
    net_listening/test_1.json
    # Erigon2 and Erigon3 never supported this api methods
    trace_rawTransaction
    debug_getRawTransaction
    # to investigate
    engine_exchangeCapabilities/test_1.json
    engine_exchangeTransitionConfigurationV1/test_01.json
    engine_getClientVersionV1/test_1.json
    # these tests require Fix on erigon DM on repeipts domain
    # these tests requires Erigon active
    admin_nodeInfo/test_01.json
    admin_peers/test_01.json
    erigon_nodeInfo/test_1.json
    eth_coinbase/test_01.json
    eth_createAccessList/test_16.json
    eth_getTransactionByHash/test_02.json
    # Small prune issue that leads to wrong ReceiptDomain data at 16999999 (probably at every million) block: https://github.com/erigontech/erigon/issues/13050
    ots_searchTransactionsBefore/test_04.tar
    eth_getWork/test_01.json
    eth_mining/test_01.json
    eth_protocolVersion/test_1.json
    eth_submitHashrate/test_1.json
    eth_submitWork/test_1.json
    net_peerCount/test_1.json
    net_version/test_1.json
    txpool_status/test_1.json
    web3_clientVersion/test_1.json)

# Transform the array into a comma-separated string
disabled_test_list=$(IFS=,; echo "${disabled_tests[*]}")

python3 ./run_tests.py -p 8545 --continue -f --json-diff --serial -x "$disabled_test_list"

exit $?
