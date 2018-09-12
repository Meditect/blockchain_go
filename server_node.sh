#!/bin/bash

GO_PATH=/Users/qiheng/Documents/blockchain_go


export NODE_ID=1000

RUN=$GO_PATH"/bin/bcg"

$RUN createwallet
$RUN createblockchain -address 1FWM9nCeSiooW71Mk9Ahr4mz45oAo87f8L
