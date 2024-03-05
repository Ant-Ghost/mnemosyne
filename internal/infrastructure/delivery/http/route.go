package http

import (
	"net/http"
	database "zocket/mnemosyne/internal/infrastructure/database/replication"
	"zocket/mnemosyne/internal/presentator/controller"

	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
)

// Responsible for creating a MUX router and register the api routes.
var RegisterDatabaseRouter = func(raftNode *raft.Raft, fsmStore database.InternalFSM) http.Handler {

	router := mux.NewRouter()

	c := controller.NewController(raftNode, fsmStore)

	router.HandleFunc("/health", c.Health)
	router.HandleFunc("/leader", c.Leader)
	router.HandleFunc("/key/{key}", c.Get).Methods("GET")
	router.HandleFunc("/key/{key}", c.Set).Methods("POST")
	router.HandleFunc("/key/{key}", c.Update).Methods("PUT")
	router.HandleFunc("/key/{key}", c.Delete).Methods("DELETE")

	return router
}
