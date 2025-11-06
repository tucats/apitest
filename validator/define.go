package validator

import "sync"

var dictionary = make(map[string]*Item)
var dictionaryLock sync.Mutex

func Define(name string, obj any) error {
	item, err := New(obj)
	if err != nil {
		return err
	}

	dictionaryLock.Lock()
	defer dictionaryLock.Unlock()

	dictionary[name] = item

	return nil
}
