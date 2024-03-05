package application

import (
	"zocket/mnemosyne/internal/domains/entity"
)

// KeyValueStoreUsecase wraps the KeyValueStore to implement required use cases
type KeyValueStoreUsecase struct {
	store entity.KeyValueStoreInterface
}

// Making sure KeyValueStoreUsercase Implements KeyValueStoreUsecaseInterface
var _ KeyValueStoreUsecaseInterface = (*KeyValueStoreUsecase)(nil)

func Init() KeyValueStoreUsecaseInterface {
	return &KeyValueStoreUsecase{
		store: entity.NewStore(),
	}
}

// Get is responsible for fetching Data from KeyValueStore
// If the value against the given key exists, then it returns the value and true.
// If the value against the given key is absent, then it returns empty string and false.
func (storeUsecase *KeyValueStoreUsecase) Get(key string) (string, bool) {
	return storeUsecase.store.GetOneValue(key)
}

// Set is responsible for registering the `key` and the `value` in KeyValueStore
func (storeUsecase *KeyValueStoreUsecase) Set(key, value string) {
	storeUsecase.store.SetOneValue(key, value)
}

// Del is repsonsible for deleting a entry from KeyValueStore using given `key`
func (storeUsecase *KeyValueStoreUsecase) Del(key string) {
	storeUsecase.store.DeleteOneValue(key)
}

// Flush is responsible for re-initializing the store, i.e. flushing all data and make the KeyValueStore Empty
func (storeUsecase *KeyValueStoreUsecase) Flush() {
	storeUsecase.store = entity.NewStore()
}

// All is reponsible for fetching the keys and values from the KeyStore  as a map
func (storeUsecase *KeyValueStoreUsecase) All() map[string]string {
	return storeUsecase.store.GetAllKeyAndValues()
}
