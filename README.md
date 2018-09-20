# blockchain tracking serial numbers

WARNING:
- hardcoded GO_PATH in scripts, which need be changed
- need more robust input error checking
- should have directories named "db" and "wallet" in the GO_PATH directory

#Summary
Coming soon

# Server node script:
```bash
# script Node_IP Node_Port API_IP API_Port
> ./server.sh localhost 3000 localhost 2000
```

GET API on server:
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


