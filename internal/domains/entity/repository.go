package entity

// Making sure KeyValueStore Implements KeyValueStoreInterface
var _ KeyValueStoreInterface = (*KeyValueStore)(nil)

// GetOneValue accepts a `key` and finds `value` associated with it in KeyValueStore.
// If the value exists, then it returns the value and true.
// If the value is absent, then empty string and false.
func (store *KeyValueStore) GetOneValue(key string) (string, bool) {
	value, exists := store.keyValue[key]
	return value, exists
}

// GetAllKeyAndValue responsible for returning the whole map of KeyValueStore
func (store *KeyValueStore) GetAllKeyAndValues() map[string]string {
	return store.keyValue
}

// SetOneValue is responsible for registering the key and its value in KeyValueStore
func (store *KeyValueStore) SetOneValue(key, value string) {
	store.keyValue[key] = value
}

// DeleteOneValue accepts one `key` and delete that value from KeyValueStore
func (store *KeyValueStore) DeleteOneValue(key string) {
	delete(store.keyValue, key)
}
