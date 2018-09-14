#!/bin/bash

export GO_PATH=/Users/qiheng/Documents/blockchain_go/
export NODE_ID=$1
echo "Running on NODE_ID: " $NODE_ID
cd $GO_PATH
DB_FILE=db/blockchain_${NODE_ID}.db
WALLET_FILE=wallet/wallet_${NODE_ID}.dat

RUN=$GO_PATH"/bin/bcg"

if ! [ -e $DB_FILE ]; then
	addr="$($RUN createwallet)"
	addr=${addr#Your new address: } #changing fmt.Printf string in source code will break this line
	echo "Your new address: " $addr
	#copy the genesis block before downloading blocks from server node
	cp db/genesis_block.db $DB_FILE
fi

echo "Addresses in wallet: "
$RUN listaddresses

$RUN startnode
