package main

import (
	"sync"
)

func transform(target *Target, extractChannel <-chan []map[string]interface{}, transformChannel chan<- []map[string]interface{}) {
	var wg sync.WaitGroup
	wg.Add(args.workers)
	for w := 0; w < args.workers; w++ {
		go func() {
			defer wg.Done()
			for data := range extractChannel {
				for _, nest := range target.Nest {
					err := queryNest(target.connection.db, nest, data)
					check(err)
				}
				transformChannel <- data
			}
		}()
	}
	wg.Wait()
	close(transformChannel)
}
