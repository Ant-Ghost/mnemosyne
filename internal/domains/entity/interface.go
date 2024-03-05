package entity

type KeyValueStoreInterface interface {
	GetOneValue(string) (string, bool)

	GetAllKeyAndValues() map[string]string

	SetOneValue(string, string)

	DeleteOneValue(string)
}
