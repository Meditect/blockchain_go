# forked from jeiwan/blockchain_go

WARNING:
- hardcoded GO_PATH into scripts, which needs be changed
- need more robust input error checking
- should have directories named "db" and "wallet" in the GO_PATH directory

# Server node script:
```bash
# script Node_IP Node_Port API_IP API_Port
> ./server.sh localhost 3000 localhost 2000
```

Server accepts API calls to query the blockchain:
get(serialnumber, salt) -> txid, PubKeyFrom, PubKeyHash
```bash
> curl --header "Content-Type: application/json" --request POST --data '{"serialnumber":"1234567890","salt":"salt"}' http://localhost:2000/get

{"Txid":null,"PubKeyFrom":null,"PubKeyHash":"QGApgHsaRW1opYwlZ15NBm0UYSw="}
```

# Client node script:
```bash
# script Node_IP Node_Port
> ./wallet.sh localhost 3001
```

The script takes client into a command line interface to:

- introduce a new serial number into blockchain
```bash
# add to_addr serialnumber salt
> add 
```

- update a serial number
```bash
# send from_addr to_addr serialnumber salt
> send 
```

- get info about a serial number
```bash
# get serialnumber salt
> 
```

# Miner node script:
```bash
# script Node_IP Node_Port Miner_PubKey
> ./miner.sh localhost 2999 xxxxxxxxx
```




Add serial number : coinbase transaction
Get serial number : GetSerialNumber(address) -> []SerialNumbers
Trace serial number : SerialNumber -> []address
Update location : transaction


















