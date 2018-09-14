#!/bin/bash

#This should get the server node up running.
export GO_PATH=/Users/qiheng/Documents/blockchain_go/
export NODE_ID=$1
cd $GO_PATH
DB_FILE=db/blockchain_${NODE_ID}.db
WALLET_FILE=wallet/wallet_${NODE_ID}.dat
pwd

RUN=bin/bcg

if ! [ -e $DB_FILE ]; then
	echo "Init server db once forever."
	addr=$($RUN createwallet)
	addr=${addr#Your new address: } #changing fmt.Printf string in source code will break this line
	$RUN createblockchain -address $addr
fi

echo "Address: " $($RUN listaddresses)

$RUN startnode
