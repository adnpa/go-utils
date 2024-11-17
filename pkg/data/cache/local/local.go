package local

type Cache interface {
	Set(key string, data []byte)
	Get(key string) ([]byte, bool)
	Del(key string)
}
