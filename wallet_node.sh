#!/bin/bash

# Set env variables
GO_PATH=/Users/qiheng/Documents/blockchain_go

# Set NODE_ID
export NODE_ID=1001


RUN=$GO_PATH"/bin/bcg"

$RUN createblockchain
