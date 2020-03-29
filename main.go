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
	go sds.checkForDeadKeys()
	return sds
}

// Set inserts {value} into the data store with an association to {key}. This
// mapping will only exist for {lifespan} seconds. After which, any subsequent
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
	if !sds.Has(key) {
		return nil
	}
	return sds.mutexGet(key).value
}

// Has returns a boolean based on whether or not the store contains a value for
// {key}.
func (sds SdrieDataStore) Has(key string) bool {
	sds.mutex.RLock()
	_, ok := sds.data[key]
	sds.mutex.RUnlock()
	return ok
}

//

func (sds SdrieDataStore) mutexGet(key string) sdrieMapValue {
	sds.mutex.RLock()
	smv := sds.data[key]
	sds.mutex.RUnlock()
	return smv
}

func (sds SdrieDataStore) mutexSet(key string, value sdrieMapValue) {
	sds.mutex.Lock()
	value.death += time.Now().Unix()
	sds.data[key] = value
	sds.mutex.Unlock()
}

func (sds SdrieDataStore) mutexDelete(key string) {
	sds.mutex.Lock()
	delete(sds.data, key)
	sds.mutex.Unlock()
}

//

func (sds SdrieDataStore) checkForDeadKeys() {
	for true {
		now := time.Now().Unix()
		toRemove := []*list.Element{}
		for e := sds.line.Front(); e != nil; e = e.Next() {
			k := e.Value.(string)
			if sds.data[k].death <= now {
				toRemove = append(toRemove, e)
				sds.mutexDelete(k)
			}
		}
		for _, item := range toRemove {
			sds.line.Remove(item)
		}
		time.Sleep(time.Second)
	}
}
