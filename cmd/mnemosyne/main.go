package main

import (
	database "zocket/mnemosyne/internal/infrastructure/database/replication"

	"zocket/mnemosyne/internal/infrastructure/delivery"
)

func main() {

	// Creating RAFT Node and FSM Store Objects
	raftNode, fsmStore := database.Replicate()

	// Initiating Delivery Protocol
	delivery.Deliver(raftNode, fsmStore)
}
