package store

import (
	"sync"

	"github.com/lewgun/argyroneta/pkg/types"
)

const (
	Bolt = "boltdb"
)

type BlobStore interface {
	AddBlob(r *types.Blob) error
	DeleteBlob(key string) error
	Blob(key string) (types.Blob, error)
}

//Store 为所有存储设备提供一个基本接口
type Store interface {
	BlobStore

	PowerOn(string) error
	PowerOff() error
}

var (
	stores = map[string]Store{}
	mu     sync.Mutex
)

//Register 注册一个存储设备
func Register(name string, store Store) {
	mu.Lock()
	defer mu.Unlock()
	stores[name] = store
}

//PowerOn 打开一个存储设备
func PowerOn(name string, conf string) (Store, error) {
	if store, ok := stores[name]; !ok {
		panic("store not found")
	} else {
		return store, store.PowerOn(conf)
	}
}
