package cache

type TimerCache interface {
	Flush(key string)
	Get(hk string, k string, wheel int) (string, error)
}
