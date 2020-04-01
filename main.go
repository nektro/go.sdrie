package sdrie

import (
	"container/list"
	"sync"
	"time"
)

type SdrieDataStore struct {
	data  map[string]sdrieMapValue
	line  *list.List
	mutex sync.RWMutex
}

type sdrieMapValue struct {
	death int64
	value interface{}
}

func New() SdrieDataStore {
	sds := SdrieDataStore{
		map[string]sdrieMapValue{},
		list.New(),
		sync.RWMutex{},
	}
	return sds
}

func (sds SdrieDataStore) Delete(key string) {
	sds.mutexDelete(key)
}

// Set inserts {value} into the data store with an association to {key}. This
// mapping will only exist for {lifespan} milliseconds. After which, any subsequent
// calls to Get will return nil unless a new value is Set.
func (sds SdrieDataStore) Set(key string, value interface{}, lifespan int64) {
	for e := sds.line.Front(); e != nil; e = e.Next() {
		k := e.Value.(string)
		if k == key {
			sds.line.Remove(e)
		}
	}
	temp := sdrieMapValue{
		lifespan,
		value,
	}
	sds.mutexSet(key, temp)
	sds.line.PushBack(key)
}

// Get retrieves the current live value associated to {key} in the store.
func (sds SdrieDataStore) Get(key string) interface{} {
	if sds.Has(key) {
		return sds.mutexGet(key).value
	} else {
		return nil
	}
}

// Has returns a boolean based on whether or not the store contains a value for {key}.
func (sds SdrieDataStore) Has(key string) bool {
	return sds.mutexHas(key)
}

//

func (sds SdrieDataStore) mutexHas(key string) bool {
	sds.mutex.RLock()
	smv, ok := sds.data[key]
	if ok && smv.death <= (time.Now().Unix() * 1000) {
		sds.unsafeDelete(key)
		ok = false
	}
	sds.mutex.RUnlock()
	return ok
}

func (sds SdrieDataStore) mutexGet(key string) sdrieMapValue {
	sds.mutex.RLock()
	smv, _ := sds.data[key]
	sds.mutex.RUnlock()
	return smv
}

func (sds SdrieDataStore) mutexSet(key string, value sdrieMapValue) {
	sds.mutex.Lock()
	value.death += time.Now().Unix() * 1000
	sds.data[key] = value
	sds.mutex.Unlock()
}

func (sds SdrieDataStore) mutexDelete(key string) {
	sds.mutex.Lock()
	sds.unsafeDelete(key)
	sds.mutex.Unlock()
}

func (sds SdrieDataStore) unsafeDelete(key string) {
	delete(sds.data, key)
}
