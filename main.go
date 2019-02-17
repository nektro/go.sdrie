package sdrie

import (
	"errors"
	"fmt"
	"time"
)

type sdrieMapValue struct {
	added    int64
	lifespan int64
	value    interface{}
}

var (
	data = map[string]sdrieMapValue{}
)

// Set adds {value} to the data store associated to {key}. If the {key} already
// exists in the store then the previous value will be overwritten. After {lifespan}
// attempting to Get {key} will return an error until a new value is added.
func Set(key string, value interface{}, lifespan int64) {
	data[key] = sdrieMapValue{
		time.Now().Unix(),
		lifespan,
		value,
	}
	death := lifespan * int64(time.Second.Seconds())
	time.AfterFunc(time.Duration(death), func() {
		if !Has(key) {
			return
		}
		if time.Now().Unix() < data[key].lifespan {
			return
		}
		delete(data, key)
	})
}

// Has returns a boolean based on whether the store contains a value for {key}
func Has(key string) bool {
	_, ok := data[key]
	return ok
}

// Get returns the value in the store associated to {key}, otherwise and nil
// if no value is found.
func Get(key string) interface{} {
	if !Has(key) {
		return nil
	}
	return data[key].value
}

// Remove preemptively removes {key} from the store even if the original
// lifespan has not passed.
func Remove(key string) error {
	if !Has(key) {
		return errors.New(fmt.Sprintf("Key '%s' not found in store", key))
	}
	delete(data, key)
	return nil
}
