package sdrie

import (
	"sync"
	"time"
)

var defaultCleanupTriggerThresold int = 1000

type SdrieDataStore struct {
	data                    map[string]sdrieMapValue
	cleanupTriggerThreshold int
	mutex                   sync.RWMutex
}

type sdrieMapValue struct {
	death int64
	value interface{}
}

func New(cleanupTriggerThreshold int) SdrieDataStore {
	if cleanupTriggerThreshold <= 0 {
		// if map exceeds threshold size, run cleanup
		cleanupTriggerThreshold = defaultCleanupTriggerThresold
	}
	sds := SdrieDataStore{
		map[string]sdrieMapValue{},
		cleanupTriggerThreshold,
		sync.RWMutex{},
	}
	return sds
}

func (sds *SdrieDataStore) Delete(key string) {
	if sds.mutexGetSize() >= sds.cleanupTriggerThreshold {
		// if map exceeds threshold size, run cleanup
		sds.mutexCleanup()
	}
	sds.mutexDelete(key)
}

// Set inserts {value} into the data store with an association to {key}. This
// mapping will only exist for {lifespan} milliseconds. After which, any subsequent
// calls to Get will return nil unless a new value is Set.
func (sds *SdrieDataStore) Set(key string, value interface{}, lifespan int64) {
	if sds.mutexGetSize() >= sds.cleanupTriggerThreshold {
		// if map exceeds threshold size, run cleanup
		sds.mutexCleanup()
	}
	temp := sdrieMapValue{
		lifespan,
		value,
	}
	sds.mutexSet(key, temp)
}

// Get retrieves the current live value associated to {key} in the store.
func (sds *SdrieDataStore) Get(key string) interface{} {
	if sds.Has(key) {
		return sds.mutexGet(key).value
	} else {
		return nil
	}
}

// Has returns a boolean based on whether or not the store contains a value for {key}.
func (sds *SdrieDataStore) Has(key string) bool {
	if sds.mutexGetSize() >= sds.cleanupTriggerThreshold {
		sds.mutexCleanup()
	}
	return sds.mutexHas(key)
}

//

func (sds *SdrieDataStore) mutexHas(key string) bool {
	sds.mutex.RLock()
	smv, ok := sds.data[key]
	if ok && smv.death <= (time.Now().Unix()*1000) {
		sds.unsafeDelete(key)
		ok = false
	}
	sds.mutex.RUnlock()
	return ok
}

func (sds *SdrieDataStore) mutexGet(key string) sdrieMapValue {
	sds.mutex.RLock()
	smv, _ := sds.data[key]
	sds.mutex.RUnlock()
	return smv
}

func (sds *SdrieDataStore) mutexSet(key string, value sdrieMapValue) {
	sds.mutex.Lock()
	value.death += time.Now().Unix() * 1000
	sds.data[key] = value
	sds.mutex.Unlock()
}

func (sds *SdrieDataStore) mutexDelete(key string) {
	sds.mutex.Lock()
	sds.unsafeDelete(key)
	sds.mutex.Unlock()
}

func (sds *SdrieDataStore) mutexGetSize() int {
	sds.mutex.Lock()
	size := sds.unsafeGetSize()
	sds.mutex.Unlock()
	return size
}

func (sds *SdrieDataStore) mutexCleanup() {
	sds.mutex.Lock()
	sds.unsafeCleanup()
	sds.mutex.Unlock()
}

func (sds *SdrieDataStore) unsafeGetSize() int {
	return len(sds.data)
}

func (sds *SdrieDataStore) unsafeCleanup() {
	// find expired keys
	keysToDelete := []string{}
	now := time.Now().Unix() * 1000
	for key, smv := range sds.data {
		if smv.death <= now {
			keysToDelete = append(keysToDelete, key)
		}
	}
	// delete expired keys
	for _, key := range keysToDelete {
		sds.unsafeDelete(key)
	}
}

func (sds *SdrieDataStore) unsafeDelete(key string) {
	delete(sds.data, key)
}
