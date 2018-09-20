# blockchain tracking serial numbers

WARNING:
- Hardcoded GO_PATH in scripts. Change it before running.

#Setup
From GO_PATH directory, run:
```bash
go get github.com/qihengchen/bcg/...
go get github.com/boltdb/bolt/...
go get -u golang.org/x/crypto/...
go install github.com/qihengchen/bcg
```

# Server node script:
```bash
# script Node_IP Node_Port API_IP API_Port
> ./server.sh localhost 3000 localhost 2000
```

GET API on server:
get(serialnumber, salt) -> txid, PubKeyFrom, PubKeyHash
```bash
> curl --header "Content-Type: application/json" --request POST --data '{"serialnumber":"1234567890","salt":"salt"}' http://localhost:2000/get
#if found, return latest transaction id and hash(PubKey of serial number owner)
{"Txid":"91f3c965d9f9c52d96f2c52549a97731a951e0e59895aeea0913af4621f60263","PubKeyHash":"23fc6fa8404aa90fdef53b34f705e1e3f350deae"}
#if not found, return empty string
{"Txid":"","PubKeyHash":""}
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
> add 137kzApkdZDFyZZtWfqwZrx2i5UU1QFtgp 1234567890 1234
Success! or error message
```
Generate serial number 1234567890 with salt 1234 and assign to a party who owns pubkey 137kzApkdZDFyZZtWfqwZrx2i5UU1QFtgp

- update a serial number
```bash
# send from_addr to_addr serialnumber salt
> send 137kzApkdZDFyZZtWfqwZrx2i5UU1QFtgp 15x5wsioY6ainZ83hQu1RMLWBew1xZF7g5 1234567890 1234
Success! or error message
```
Send serial number 1234567890 with salt 1234 from 137kzApkdZDFyZZtWfqwZrx2i5UU1QFtgp to another address 15x5wsioY6ainZ83hQu1RMLWBew1xZF7g5

- get info about a serial number
```bash
# get serialnumber salt
> get 1234567890 1234
============ Tx %!x(MISSING) ============
Serial Number Hash: [194 96 9 245 130 251 24 78 86 1 20 118 194 251 161 97 55 41 215 83 82 87 19 68 150 47 201 41 182 225 82 177]
Prev. block: 364b8284ddd8ca5caa3365326b80c8c6531fe866
```

- print the blockchain
```bash
# print the whole blockchain
# TODO: print the latest n blocks
> print
```

# Miner node script:
```bash
# script Node_IP Node_Port Miner_PubKey
> ./miner.sh localhost 2999
```


