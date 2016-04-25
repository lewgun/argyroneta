package mysqlmgr

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"Pay-Platform/pkg/errutil"
	"Pay-Platform/pkg/funcall"
	"Pay-Platform/pkg/store"
	"Pay-Platform/pkg/types"
)

//StoreMgr 多个store实例管理器
type MysqlMgr struct {
	stores   map[string]*store.Store   //store实例与其名称的映射
	doBy     map[string][]*store.Store //某个action对应的执行store(可以多于一个,由谁执行此动作则随机)
	caller   *funcall.Caller
	typeName string
}

var storeMgr *MysqlMgr

func init() {
	storeMgr = &MysqlMgr{
		stores: make(map[string]*store.Store),
		doBy:   make(map[string][]*store.Store),
		caller: funcall.New(),
	}

	st := &store.Store{}
	t := reflect.TypeOf(st)

	storeMgr.typeName = t.Elem().Name()

	storeMgr.caller.Register(st)
}

//BootUp boot up the store mgr.
func BootUp(opts map[string]*types.StoreOpt) error {

	if opts == nil {
		return errutil.ErrInvalidParameter
	}

	for k, v := range opts {
		st := store.New(v.DSN, v.ShowSQL, v.ChanLen)
		if st == nil {
			return fmt.Errorf("New store: %s's instance is failed.", k)
		}

		//缓存store
		storeMgr.stores[k] = st

		//建立store/handler与action之间的映射
		for _, act := range v.Actions {
			if _, ok := storeMgr.doBy[act]; !ok {
				storeMgr.doBy[act] = make([]*store.Store, 0)
			}
			storeMgr.doBy[act] = append(storeMgr.doBy[act], st)

		}

	}

	return nil
}

func powerOffHelper() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for _, st := range storeMgr.stores {
			wg.Add(1)
			go func(s *store.Store) {
				defer wg.Done()
				s.PowerOff()
			}(st)
		}

		wg.Wait()
		ch <- struct{}{}

	}()

	return ch
}

func PowerOff() error {

	chanExit := powerOffHelper()

	var err error
	select {
	case <-chanExit:
	case <-time.After(10 * 1e9):
		err = fmt.Errorf("timeout")

	}

	return err

}

func Do(action string, params ...interface{}) (interface{}, error) {

	st, ok := storeMgr.doBy[action]
	if !ok {
		return nil, fmt.Errorf("no store for action: %s", action)
	}

	retVal, err := storeMgr.caller.CallOnObject(st[0], storeMgr.typeName+"."+action, params...)
	if err != nil {
		return nil, err
	}

	if retVal[1] != nil {
		err = retVal[1].(error)
	}

	return retVal[0], err

}
