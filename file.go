package main

import (
	"io/ioutil"
	"sync"
)

var cache = make(map[string]string, 0)
var mutex = new(sync.Mutex)

func loadFile(name string) (string, error) {
	if _, ok := cache[name]; !ok {
		mutex.Lock()
		buffer, err := ioutil.ReadFile(name)
		if err != nil {
			return "", err
		}
		cache[name] = string(buffer)
		mutex.Unlock()
	}
	return cache[name], nil
}