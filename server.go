package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"encoding/json"
	"bufio"
	"os"
	"strings"
	"strconv"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000"} //knownNodes will be hardcoded 
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type verzion struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

type getJSONReq struct {
	SerialNumber string `json:"serialnumber"`
	Salt 		 string `json:"salt"`
}

type getJSONResp struct {
	Txid 		string
	PubKeyFrom 	string
	PubKeyHash 	string
}

//{"Txid":null,"PubKeyFrom":null,"PubKeyHash":"QGApgHsaRW1opYwlZ15NBm0UYSw="}
//{"Txid":"","PubKeyFrom":"","PubKeyHash":"406029807b1a456d68a58c25675e4d066d14612c"}

/* 	client node: address, "", ""
	miner node:	 address, "", apiAddress
	server node: address, minerAddress, "" where address = minerAddress
*/

func StartServer(address string, minerAddress, apiAddress string) {
	nodeID := ParseNodeID(address)
	nodeAddress = address
	miningAddress = minerAddress
	bc := NewBlockchain(nodeID)

	var wg sync.WaitGroup
	wg.Add(1)
	go launchTCPListener(nodeAddress, bc) //node communication

	// if server, launch API listener
	if SliceContainsString(knownNodes, nodeAddress) {
		wg.Add(1)
		go launchHTTPListener(apiAddress, bc)
	} else {
		//if not, check in with server node. Default to the first server
		sendVersion(knownNodes[0], bc) 
	}

	if minerAddress == "" && apiAddress == "" {
		launchClientInterface(bc)
	}

	wg.Wait()
}

func launchClientInterface(bc *Blockchain) {
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		mineNow = false

		line := strings.TrimSuffix(string(sentence), "\n")
		args := strings.Split(line, " ")

		//TODO: error checking
		fmt.Printf("Received %s command\n", args[0])

		switch args[0] {
		case "add":
			to := args[1]
			data := args[2]
			salt := args[3]
			clientAddHandler(to, data, salt, bc, mineNow)
			
		case "get":
			data := args[1]
			salt := args[2]
			clientGetHandler(data, salt, bc)
		case "send":
			from := args[1]
			to := args[2]
			data := args[3]
			salt := args[4]
			clientSendHandler(from, to, data, salt, bc, mineNow)
		case "print":
			clientPrintHandler(bc)
		default:
			fmt.Println("Unknown command!")
		}
	}
}

func launchTCPListener(nodeAddress string, bc *Blockchain) {
	listener, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}


func launchHTTPListener(apiAddress string, bc *Blockchain) {
	//serial number, salt -> txid, from pubkey, to pubkey hash
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		req := getJSONReq{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != io.EOF && err != nil {
			panic(err)
		}
		fmt.Println(req)
		resp := getJSONResp{}

		hash := HashSerialNumber(req.SerialNumber, req.Salt)
		outputs := bc.FindSerialNumberHash(hash)
		if len(outputs) > 0 {
			resp.PubKeyHash = fmt.Sprintf("%x", outputs[0].PubKeyHash)
		}
		fmt.Println(resp)
		respJson, err := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	})

	log.Fatal(http.ListenAndServe(apiAddress, nil))
	
	/*s := &http.Server{
		Addr:           apiAddress
		Handler:        handler_name,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())*/
}

func clientAddHandler(to string, serialNumber, salt string, bc *Blockchain, mineNow bool) {
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	tx := NewSerialNumberTX(to, serialNumber, salt)

	if mineNow {
		UTXOSet := UTXOSet{bc}
		txs := []*Transaction{tx}
		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")
}

func clientGetHandler(serialNumber, salt string, bc *Blockchain) {
	UTXOSet := UTXOSet{bc}

	hash := HashSerialNumber(serialNumber, salt)

	outputs := UTXOSet.FindSerialNumberHash(hash)
	
	for _, output := range outputs {
		fmt.Printf("============ Tx %x ============\n")
		fmt.Printf("Serial Number Hash: %d\n", output.SerialNumberHash)
		fmt.Printf("Prev. block: %x\n", output.PubKeyHash)
		fmt.Printf("\n")
	}
}

func clientSendHandler(from, to string, serialNumber, salt string, bc *Blockchain, mineNow bool) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	UTXOSet := UTXOSet{bc}

	wallets, err := NewWallets(ParseNodeID(nodeAddress))
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)

	// later may be modified to transfer labels in patch
	tx, err := NewUTXOTransaction(&wallet, to, serialNumber, salt, &UTXOSet)

	if err != nil {
		fmt.Println(err)
		return
	}

	if mineNow {
		txs := []*Transaction{tx}
		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")
}

func clientPrintHandler(bc *Blockchain) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

// Send data to addr
func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}

// TODO: validate blocks
func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

// TODO: check we actually have the block or tx
func handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			//TODO: salt parameter
			cbTx := NewSerialNumberTX(miningAddress, "", "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := UTXOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload verzion

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func handleConnection(conn net.Conn, bc *Blockchain) {

	request, err := ioutil.ReadAll(conn)
	//request := make([]byte, 2048)
	//_, err := conn.Read(request)
	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}
	conn.Close()
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
