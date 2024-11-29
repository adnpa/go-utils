package local

type Local interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	Del(key string)
}
