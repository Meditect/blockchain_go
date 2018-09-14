#!/bin/bash

export GO_PATH=/Users/qiheng/Documents/blockchain_go
export NODE_ID=$1
echo "NODE_ID: " $NODE_ID
cd $GO_PATH
DB_FILE=db/blockchain_${NODE_ID}.db
WALLET_FILE=wallet/wallet_${NODE_ID}.dat

RUN=$GO_PATH"/bin/bcg"


if ! [ -e DB_FILE ]; then
	addr="$($RUN createwallet)"
	addr=${addr#Your new address: }
	echo "Your new address: " $addr
	cp db/blockchain_3000.db $DB_FILE
fi

$RUN startnode -miner $addr
