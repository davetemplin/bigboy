package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

var errors uint64

func main() {
	var progname string = os.Args[0]
	var input []string = os.Args[1:]

	flags, usage := parseFlags(progname, input, &config)
	if flags.version {
		stop(version, 0)
	}

	config = *loadConfig(flags.config)
	args = *parseArgs(progname, input, &config)

	if args.target == "" {
		usage()
		os.Exit(0)
	}

	t := time.Now()
	target := loadTarget()
	run(target)
	disconnect()
	if !args.quiet {
		fmt.Printf("%s: %d seconds elapsed\n", args.target, int(time.Since(t).Seconds()))
	}
}

func run(target *Target) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(target *Target) {
		extractChannel := make(chan []map[string]interface{}, args.page*args.workers*10)
		transformChannel := make(chan []map[string]interface{}, args.page*args.workers*10)
		go extract(target, extractChannel)
		go transform(target, extractChannel, transformChannel)
		go write(target, transformChannel, &wg)
	}(target)
	wg.Wait()
}
