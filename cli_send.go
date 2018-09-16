package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from string, to, serialNumber string, salt, nodeID string, mineNow bool) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)

	// later may be modified to transfer labels in patch
	// func NewUTXOTransaction(wallet *Wallet, addresses, serialNumbers []string, 
	//             salt string, UTXOSet *UTXOSet) (*Transaction, []string)
	//to_array := []string{"a", "b"}
	//serialNumber_array := []string{"c", "d"}
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

/*
func jsonToArray(json string, flag, string) []string {
	if flag == "to" {

	}
	if flag == "serialNumber" {

	}
}*/


