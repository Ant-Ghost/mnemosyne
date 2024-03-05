package database

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"sync"

	"zocket/mnemosyne/internal/application"

	"github.com/hashicorp/raft"
)

// KeyValueFSM is used to Implement the Finite State Machine.
// It is responsible for the state of the ServerNode
// It contains a Read-Write Mutex field `mu` for safe reading and writing data in KeyValueStore
// It contains a KeyValueStoreUsecase field `store` to perform multiple operations on KeyValueStore
type KeyValueFSM struct {
	mu           sync.RWMutex
	storeUsecase application.KeyValueStoreUsecaseInterface
}

// KeyValueCommand is responsible for enforcing the format for LogEntries from same/different server node
type KeyValueCommand struct {
	Operation string `json:"operation"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// Making sure that KeyValueFSM implements InternalFSM
var _ InternalFSM = (*KeyValueFSM)(nil)

// NewFSM is responsible for initializing storeUsecase and creating a new KeyValueFSM, i.e. a new state of the Machine
func NewFSM() InternalFSM {
	usecase := application.Init()
	return &KeyValueFSM{storeUsecase: usecase}
}

// Restore is responsible writing the given Snapshot in the KeyValueStore of the Machine
func (fsm *KeyValueFSM) Restore(snapshot io.ReadCloser) error {

	// Fetching the data in a byte slice
	data, err := io.ReadAll(snapshot)
	if err != nil {
		return err
	}

	// Converting the byte slice to a string to string map
	var snapshotData map[string]string
	if err := json.Unmarshal(data, &snapshotData); err != nil {
		return err
	}

	// Implementing a Write Lock
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	// Flushing the whole keystore and writing the Snapshot data in KeyValueStore
	fsm.storeUsecase.Flush()
	for key, value := range snapshotData {
		fsm.storeUsecase.Set(key, value)
	}

	return nil
}

// Snapshot is responsible for creating a snapshot of current KeyValueStore of the Machine
func (fsm *KeyValueFSM) Snapshot() (raft.FSMSnapshot, error) {

	// Implementing a Read Lock
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()

	// Creating a snapshot variable and loading the current KeyValueStore data in it
	snapshot := make(map[string]string)
	keyValues := fsm.storeUsecase.All()
	for key, value := range keyValues {
		snapshot[key] = value
	}

	return &KeyValueSnapshot{data: snapshot}, nil
}

// Apply is responsible for applying a Write Command like Setting a Key-Value or Deleting a Key-Value on the KeyValueStore of the Machine
func (fsm *KeyValueFSM) Apply(logData *raft.Log) interface{} {

	// Implementing only for a LogCommand
	if logData.Type == raft.LogCommand {

		// Converting logData to a structured KeyValueCommand
		var command KeyValueCommand
		if err := json.Unmarshal(logData.Data, &command); err != nil {
			log.Printf("Failed to unmarshal log data: %v", err)
			return nil
		}

		// Making sure that the operations are in uppercase
		operation := strings.ToUpper(strings.TrimSpace(command.Operation))

		switch operation {
		case "SET":
			fsm.set(command.Key, command.Value)
		case "DELETE":
			fsm.delete(command.Key)
		}
	}

	return nil

}

// `Read` is responsible for fetcing value by given key from the KeyValueStore of the Machine
func (fsm *KeyValueFSM) Read(key string) (string, bool) {

	// Implementing a Read Lock
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()

	value, exists := fsm.storeUsecase.Get(key)

	return value, exists
}

// `set` sets a key-value pair in the KeyValueStore of the Machine
func (fsm *KeyValueFSM) set(key, value string) {
	// Implementing a Write Lock
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	fsm.storeUsecase.Set(key, value)
}

// `delete` deletes a key from the KeyValueStore of the Machine
func (fsm *KeyValueFSM) delete(key string) {
	// Implementing a Write Lock
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	fsm.storeUsecase.Del(key)
}
