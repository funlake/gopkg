package cache

type KvStore interface {
	Get(key string,opts ... interface{}) (interface{}, error)
	Set(key string, val interface{}) (interface{},error)
}
