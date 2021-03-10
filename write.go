package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

func write(target *Target, transformChannel <-chan []map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	if target.Split == nil {
		writeAll(target, transformChannel)
	} else {
		writeSplit(target, transformChannel)
	}
}

func writeAll(target *Target, transformChannel <-chan []map[string]interface{}) {
	t := time.Now()

	var fileName string
	if args.out == "" {
		fileName = path.Join(defaultOutputDir, args.target, day(t)+".json")
	} else {
		fileName = args.out
	}

	err := os.MkdirAll(path.Dir(fileName), 0777)
	check(err)

	file, err := os.Create(fileName)
	check(err)
	defer file.Close()

	rows := 0
	for data := range transformChannel {
		for _, obj := range data {
			err := jsonWriteln(file, obj)
			check(err)
		}
		rows += len(data)
		if !args.quiet {
			fmt.Printf("%d rows written\r", rows)
		}
	}

	if !args.quiet {
		fmt.Printf("%d rows written in %d seconds\n", rows, int(time.Since(t).Seconds()))
	}
}

func writeSplit(target *Target, transformChannel <-chan []map[string]interface{}) {
	t := time.Now()

	var dir string
	if args.out == "" {
		dir = path.Join(defaultOutputDir, args.target)
	} else {
		dir = path.Join(path.Dir(args.out), args.target)
	}

	err := os.MkdirAll(path.Dir(dir), 0777)
	check(err)

	s := newSplitter(target, dir)
	rows := 0
	for data := range transformChannel {
		for _, obj := range data {
			file := s.split(obj)
			err := jsonWriteln(file, obj)
			check(err)
		}
		rows += len(data)
		if !args.quiet {
			fmt.Printf("%d rows written, %d files\r", rows, len(s.files))
		}
	}
	s.close()

	if !args.quiet {
		fmt.Printf("%d rows written to %d files in %d seconds\n", rows, len(s.files), int(time.Since(t).Seconds()))
	}
}
