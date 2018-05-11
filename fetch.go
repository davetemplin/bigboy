package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

func prefetch(target *Target, extractChannel chan<- []map[string]interface{}) {
	fetchChannel := make(chan []uint64, args.workers * 10)
	
	fmt.Println("prefetching...")
	filename, n := prefetchToFile(target)
	defer os.Remove(filename)
	fmt.Printf("%d rows prefetched\n", n)

	var wg sync.WaitGroup
	wg.Add(args.workers)
	for w := 0; w < args.workers; w++ {
		go fetchWorker(target, fetchChannel, extractChannel, &wg)
	}

	loadFetchChannel(filename, fetchChannel)
	close(fetchChannel)

	err := os.Remove(filename)
	check(err)

	wg.Wait()
}

func fetchWorker(target *Target, fetchChannel <-chan []uint64, extractChannel chan<- []map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for ids := range fetchChannel {
		var retry uint64
		for {
			result, err := extractIds(target, ids)
			if (err == nil) {
				for _, data := range result {
					extractChannel <- data
				}
				break
			}
			retryCheck(err, &retry)
		}
	}
}

func prefetchToFile(target *Target) (string, int) {
	n := 0
	rows, err := target.connection.db.Query(target.prefetch, target.params...)
	check(err)
	defer rows.Close()

	file, err := ioutil.TempFile("", "bigboy-")
	check(err)
	defer file.Close()

	var id uint64
	bytes := make([]byte, 8)
	for rows.Next() {
		err = rows.Scan(&id)
		check(err)
		binary.LittleEndian.PutUint64(bytes, uint64(id))
		_, err := file.Write(bytes)
		check(err)
		n++
	}

	err = file.Close()
	check(err)

	return file.Name(), n
}

func loadFetchChannel(filename string, ch chan<- []uint64) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)
	bytes := make([]byte, 8)
	page := make([]uint64, 0)
	for {
		_, err := io.ReadFull(r, bytes)
		if err == io.EOF {
			if len(page) > 0 {
				ch <- page
			}
			break
		} else {
			check(err)
		}

		value := binary.LittleEndian.Uint64(bytes)
		page = append(page, value)
		if len(page) == args.page {
			ch <- page
			page = make([]uint64, 0)
		}
	}

	err = f.Close()
	check(err)	
}