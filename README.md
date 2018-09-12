# forked from jiewan/blockchain_go

## TODO:
- built the distributed network part
- start up script to run a code, for each pharma company
- API put, get, update serial numbers
- testing

## Issues:
- REQ on serial numbers
- assuming each pharma company run miner, full and client node?


Monday:
wallet.go
wallets.go
transaction.go
transaction_input.go
transaction_output.go

Tuesday:
block.go
blockchain.go
cli.go
cli_*.go
proofofwork.go
utxo_set.go

Wednesday:
server.go
startup.sh
APIs


Add serial number : coinbase transaction
Get serial number : GetSerialNumber(address) -> []SerialNumbers
Trace serial number : SerialNumber -> []address
Update location : transaction


type Transaction struct {
    id
    Input
    Output
}

type Input struct {
    reference_to_previous_output
    signature
    // hash(serial number + salt1) + hash(location1 + salt1)
    pubKey of checkpoint1
}

type Output struct {
    hash(serial number + salt2)
    hash(location2 + salt2)
    pubKey of checkpoint2
}



















