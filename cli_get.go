package main

import (
	"fmt"
)

func (cli *CLI) getSerialNumber(serialNumber, salt string, nodeID string) {
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	hash := HashSerialNumber(serialNumber, salt)

	outputs, txIDs := UTXOSet.FindSerialNumberHash(hash)
	
	for index, output := range outputs {
		fmt.Printf("============ Transaction %x ============\n", txIDs[index])
		fmt.Printf("Serial Number Hash: %d\n", output.SerialNumberHash)
		fmt.Printf("Script (PubKey hash of recipient): %x\n", output.PubKeyHash)
		fmt.Printf("\n")
	}
}
