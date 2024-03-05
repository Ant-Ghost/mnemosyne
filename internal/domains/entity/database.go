package entity

// KeyValueStore is the most important entity as it is responsible for storing the key-value data.
// KeyValueStore is a structure with field `keyValue` that stores string to string map
type KeyValueStore struct {
	keyValue map[string]string
}

// NewStore is responsible for creating a brand new KeyValueStore instance
func NewStore() KeyValueStoreInterface {
	return &KeyValueStore{keyValue: make(map[string]string)}
}
