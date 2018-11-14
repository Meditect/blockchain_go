# blockchain tracking serial numbers

WARNING:
- Hardcoded GO_PATH in scripts. Change it before running.

Setup: in GO_PATH directory, run:
```bash
go get github.com/qihengchen/blockchain_go/...
go get github.com/boltdb/bolt/...
go get -u golang.org/x/crypto/...
# used 'bcg' in scripts, so renaming the directory
mv src/github.com/qihengchen/blockchain_go src/github.com/qihengchen/bcg
go install github.com/qihengchen/bcg
```


Scenario:
1. Start "server" node. All clients can add and update. This node is server-like only in the sense that miners send new blocks to the server, and the server broadcasts to all clients. We can change it later, but for now, this is good.

- ./server.sh localhost 3000 localhost 2000

2. Start miner. Note that we configured the miner to start mining after receiving two tx. In a new window, type:

- ./miner.sh localhost 2999

3. Start client one and create two labels. In a new window, type:

- ./client.sh localhost 3001

In the dummy CLI, type:
- add [paste its address from stdout] 30013001 0   (30013001 is the serial number; 0 is a placeholder for the legacy 'salt' field)
- add [paste its address from stdout] 300130013001 0

4. Start client two, create a label, and transfer an existing label. In a new window, type:

- ./client.sh localhost 3002

In the dummy CLI, type:
- add [paste its address from stdout] 30023002 0

Back to client one's window, type:
- send [node one addr] [node two addr] 30013001 0  (transferring '30013001' from node one to node two)

5. Start client three, create a label, and transfer the '30013001' label again. In a new window, type:

- ./client.sh localhost 3001

In the dummy CLI, type:
- add [paste node three address from stdout] 30033003 0
- send [node two addr] [node three addr] 30013001 0  (the '30013001' is transferred again to node three, after being transfer from node one to node two)







========= IGNORE ==========
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


# Miner node script:
```bash
# script Node_IP Node_Port Miner_PubKey
> ./miner.sh localhost 2999
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
============ Transaction ============
txid: 31643961313532623935323965373763303130613261316537343739323535303733623031663664303062646535363338396436666134343634313761356233
Serial Number Hash: [91 130 252 3 108 94 148 130 48 79 163 48 9 19 245 56 88 65 73 245 251 57 116 113 25 252 7 224 20 248 219 174]
Script (PubKey hash of recipient): 5b7a1be33dcba30823e3fe51778ba274cadb7dd5
```

- print the blockchain
```bash
# print the whole blockchain
# TODO: print the latest n blocks
> print
```




