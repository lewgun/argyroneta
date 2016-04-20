package memory

import ( 
	"sync"
    
	"github.com/lewgun/argyroneta/pkg/errutil"
)

type Memory struct {
	sync.RWMutex

	items map[string]interface{}
}

func New() *Memory {
	return &Memory{
		items: map[string]interface{}{},
	}
}

func (r *Memory) Get(key string) (interface{}, error) {
	r.RLock()
	defer r.RUnlock()

	value, ok := r.items[key]
	if !ok {
		return nil, errutil.ErrNotFound
	}

	return value, nil
}

func (r *Memory) Set(key string, value interface{}) error {
	r.Lock()
	defer r.Unlock()

	r.items[key] = value
	return nil
}

func (r *Memory) SetNX(key string, value interface{}) (bool, error) {
	r.Lock()
	defer r.Unlock()

	_, ok := r.items[key]
	if ok {
		return false, nil
	} else {
		r.items[key] = value
		return true, nil
	}
}

func (r *Memory) Delete(key string) error {
	r.Lock()
	defer r.Unlock()

	delete(r.items, key)
	return nil
}

func (r *Memory) Has (key string ) bool {
    r.Lock()
	defer r.Unlock()
    
    _, ok :=  r.items[key]
    return ok 
}