package database

import (
	"encoding/gob"

	"github.com/hashicorp/raft"
)

// KeyValueSnapshot is responsible storing the state of the Machine at a given point of time
// KeyValueSnapshot contains a `data` field that store the KeyValueStore data
type KeyValueSnapshot struct {
	data map[string]string
}

// Making sure KeyValueSnapshot Implements raft.FSMSnapshot
var _ raft.FSMSnapshot = (*KeyValueSnapshot)(nil)

// Persist saves the snapshot in the sink to make sure the data is Persistance
func (snap *KeyValueSnapshot) Persist(sink raft.SnapshotSink) error {
	// Write the snapshot data to the sink
	err := gob.NewEncoder(sink).Encode(snap.data)
	if err != nil {
		sink.Cancel()
		return err
	}
	sink.Close()
	return nil
}

// Release is invoked after all operations on the current Snapshot is over.
// No logic required
func (snap *KeyValueSnapshot) Release() {}
