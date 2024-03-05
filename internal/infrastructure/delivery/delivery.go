package delivery

import (
	database "zocket/mnemosyne/internal/infrastructure/database/replication"
	"zocket/mnemosyne/internal/infrastructure/delivery/http"

	"github.com/hashicorp/raft"
)

// Responsible for Initiating Multiple Delivery Servers
func Deliver(raftNode *raft.Raft, fsmStore database.InternalFSM) {

	// Initiating an HTTP server
	http.Server(raftNode, fsmStore)

}
