package main

import (
	"fmt"
)

func (cli *CLI) getSerialNumber(serialNumber, salt string, nodeID string) {
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	hash := HashSerialNumber(serialNumber, salt)

	outputs := UTXOSet.FindSerialNumberHash(hash)
	
	for _, output := range outputs {
		fmt.Printf("============ Tx %x ============\n")
		fmt.Printf("Serial Number Hash: %d\n", output.SerialNumberHash)
		fmt.Printf("Prev. block: %x\n", output.PubKeyHash)
		fmt.Printf("\n")
	}
}
