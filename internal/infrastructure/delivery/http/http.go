package http

import (
	"log"
	"net/http"
	"os"

	database "zocket/mnemosyne/internal/infrastructure/database/replication"

	"github.com/hashicorp/raft"
)

func Server(raftNode *raft.Raft, fsmStore database.InternalFSM) {

	// Create a new HTTP handler.
	router := RegisterDatabaseRouter(raftNode, fsmStore)

	// Start the HTTP server.
	log.Println("HTTP server is now listening on port 4000.")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), router))

}
