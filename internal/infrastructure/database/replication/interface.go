package database

import "github.com/hashicorp/raft"

// A custom FSM interface by embedding the raft.FSM model
type InternalFSM interface {
	raft.FSM

	// `Read` is responsible for fetching Data from KeyValueStore of the Machine(Server Node)
	// `Read` is available for external use because it will not affect the state Machine.
	Read(string) (string, bool)

	// `set` and `delete` are setting data and delete data from the Machine respectively.
	// `set` and `delete` is internally used by Apply() as the will affect the state of the Machine.
	set(string, string)
	delete(string)
}
