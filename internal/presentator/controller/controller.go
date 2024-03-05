package controller

import (
	"encoding/json"
	"net/http"
	"time"
	database "zocket/mnemosyne/internal/infrastructure/database/replication"

	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
)

// To enforce strict formatting while working with key-value data
type KeyValueModel struct {
	Key   string
	Value string
}

// A Controller Struct that stores the raftNode and fsmStore (state of the Machine)
type Controller struct {
	raftNode *raft.Raft
	fsmStore database.InternalFSM
}

// Making sure Controller Implements ControllerType
var _ ControllerType = (*Controller)(nil)

// Create a new Controller
func NewController(raftNode *raft.Raft, fsmStore database.InternalFSM) ControllerType {
	return &Controller{
		raftNode: raftNode,
		fsmStore: fsmStore,
	}
}

// Health Controller is responsible for signaling the the connection status of the Server
// It always returns HTTP Status OK
func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Leader Controller is responsible telling whether the current node is Leader in the Raft Cluster
// If the current Server Node is a Leader, it returns HTTP Status OK
// If the current Server Node is not a Leader, it returns HTTP Status Service Unavailable
func (c *Controller) Leader(w http.ResponseWriter, r *http.Request) {

	// Getting the Current State of the Node
	state := c.raftNode.State()

	// Checking if the current node is not a leader
	if state != raft.Leader {

		// Current Node is not a Leader
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Current Node is a Leader
	w.WriteHeader(http.StatusOK)
}

// Get Controller is used to read value by the given key from the KeyValueStore of the current Node
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	key := params["key"]

	// Reading the value by key from the Machine(Node)
	value, exists := c.fsmStore.Read(key)

	if !exists {
		w.WriteHeader(http.StatusNotFound)
	}

	keyValue := KeyValueModel{Key: key, Value: value}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keyValue)
}

// Set Controller is responsible for initiating "SET" command to current Machine(Node) upon receiving set key-value request
func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Value string
	}

	var value requestBody
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := mux.Vars(r)["key"]

	// Creating "SET" command
	command := database.KeyValueCommand{
		Operation: "SET",
		Key:       key,
		Value:     value.Value,
	}

	// Converting the commands to byte slice
	commandBytes, err := json.Marshal(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Confirming that the current Node is a Leader before applying a Write operation
	if c.raftNode.State() != raft.Leader {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Applying the command with timeout of 500 milliseconds
	if err := c.raftNode.Apply(commandBytes, 500*time.Millisecond).Error(); err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update Controller is responsible for initiating "SET" command to current Machine(Node) upon receiving update value by key request
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)["key"]
	value := r.URL.Query().Get("value")

	// Creating "SET" command
	command := database.KeyValueCommand{
		Operation: "SET",
		Key:       key,
		Value:     value,
	}

	// Converting the commands to byte slice
	commandBytes, err := json.Marshal(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Confirming that the current Node is a Leader before applying a Write operation
	if c.raftNode.State() != raft.Leader {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Applying the command with timeout of 500 milliseconds
	if err := c.raftNode.Apply(commandBytes, 500*time.Millisecond).Error(); err != nil {
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Delete Controller is responsible for initiating "DELETE" command to current Machine(Node) upon receiving set key-value request
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	// Creating "DELETE" command
	command := database.KeyValueCommand{
		Operation: "DELETE",
		Key:       key,
		Value:     "",
	}

	// Converting the commands to byte slice
	commandBytes, err := json.Marshal(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Confirming that the current Node is a Leader before applying a Write operation
	if c.raftNode.State() != raft.Leader {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Applying the command with timeout of 500 milliseconds
	if err := c.raftNode.Apply(commandBytes, 500*time.Millisecond).Error(); err != nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
