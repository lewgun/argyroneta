package store

import (
	"sync"

	"github.com/lewgun/argyroneta/pkg/rule"
	"github.com/lewgun/argyroneta/pkg/types"
)

const (

	Bolt = "boltdb"
)

type RuleStore interface {
	AddRule(r *rule.Rule) error
	UpdateRule(r *rule.Rule) error
	DeleteRule(key string) error
	Rule(key string) (*rule.Rule, error)
	Rules() (map[string]*rule.Rule, error)
}

type BlobStore interface {
	AddBlob(r *types.Blob) error
	DeleteBlob(key string) error
	Blob(key string) (*rule.Rule, error)
}

//Store 为所有存储设备提供一个基本接口
type Store interface {
	RuleStore
	BlobStore

	Connect(string) error
	Close() error
}

var (
	stores   = map[string]Store{}
	storesMu sync.Mutex
)

//Register 注册一个存储设备
func Register(name string, store Store) {
	storesMu.Lock()
	defer storesMu.Unlock()
	stores[name] = store
}

//Open 打开一个存储设备
func Open(name string, conf string) (Store, error) {
	if store, ok := stores[name]; !ok {
		panic("store not found")
	} else {
		return store, store.Connect(conf)
	}
}
