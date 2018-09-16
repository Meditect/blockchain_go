package main

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

serverID := "3000"
node1 := "3001"
node2 := "3002"
node3 := "3003"
node4 := "3004"

func TestScenario(t *testing.T) {
	cli := CLI{}
	cli.Run()
	
	cli.startNode(serverID, "")

}