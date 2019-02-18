package sdrie

import (
	"container/list"
	"sync"
	"time"
)

type sdrieMapValue struct {
	death int64
	value interface{}
}

var (
	data  = map[string]sdrieMapValue{}
	line  = list.New()
	mutex = sync.RWMutex{}
)

// Set inserts {value} into the data store with an association to {key}. This
// mapping will only exist for {lifespan} seconds. After which, any subsequent
// calls to Get will return nil unless a new value is Set.
func Set(key string, value string, lifespan int64) {
	for e := line.Front(); e != nil; e = e.Next() {
		k := e.Value.(string)
		if k == key {
			line.Remove(e)
		}
	}
	temp := sdrieMapValue{
		time.Now().Unix() + lifespan,
		value,
	}
	mutexSet(key, temp)
	line.PushBack(key)
}

// Get retrieves the current live value associated to {key} in the store.
func Get(key string) interface{} {
	if !Has(key) {
		return nil
	}
	return mutexGet(key).value
}

// Has returns a boolean based on whether or not the store contains a value for
// {key}.
func Has(key string) bool {
	mutex.RLock()
	_, ok := data[key]
	mutex.RUnlock()
	return ok
}

//

func mutexGet(key string) sdrieMapValue {
	mutex.RLock()
	smv := data[key]
	mutex.RUnlock()
	return smv
}

func mutexSet(key string, value sdrieMapValue) {
	mutex.Lock()
	data[key] = value
	mutex.Unlock()
}

func mutexDelete(key string) {
	mutex.Lock()
	delete(data, key)
	mutex.Unlock()
}

//

func init() {
	go checkForDeadKeys()
}

func checkForDeadKeys() {
	for true {
		now := time.Now().Unix()
		toRemove := []*list.Element{}
		for e := line.Front(); e != nil; e = e.Next() {
			k := e.Value.(string)
			if data[k].death <= now {
				toRemove = append(toRemove, e)
				mutexDelete(k)
			}
		}
		for _, item := range toRemove {
			line.Remove(item)
		}
		time.Sleep(time.Second)
	}
}
