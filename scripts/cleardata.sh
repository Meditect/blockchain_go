#!/bin/bash

export GO_PATH=/Users/qiheng/Documents/blockchain_go
cd $GO_PATH
rm -rf db/blockchain_*
rm -rf wallet/wallet_*
echo "All data cleared except genesis_block.db"
