package cache


type Cache interface {
	Get(key string) (interface{}, error)

	Set(key string, value interface{}) error

	SetNX(key string, value interface{}) (bool, error)

	Delete(key string) error
}

