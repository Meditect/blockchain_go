package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, serialNumber, salt string, nodeID string, mineNow bool) {
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

	// func NewUTXOTransaction(wallet *Wallet, addresses, serialNumbers []string, 
	//             salt string, UTXOSet *UTXOSet) (*Transaction, []string)
	tx, _ := NewUTXOTransaction(&wallet, to, serialNumber, salt, &UTXOSet)

	if mineNow {
		cbTx := NewSerialNumberTX(from, "")
		txs := []*Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")
}
