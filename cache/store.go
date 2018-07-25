package cache

type KvStore interface{
	Connect(dsn string,pwd string)
	Get(key string) (interface{},error)
	Set(key string,val interface{})
	GetPool() interface{}
}