package cache

type Store interface{
	Connect()
	Get(key string)
	Set(key string,val interface{})
}