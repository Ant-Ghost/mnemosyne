package application

type KeyValueStoreUsecaseInterface interface {
	Get(string) (string, bool)
	Set(string, string)
	Del(string)
	Flush()
	All() map[string]string
}
