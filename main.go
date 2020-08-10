package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

var args struct {
	errors  uint64
	nulls   bool
	out     string
	page    int
	params  []string
	path    string
	quiet   bool
	retries uint64
	target  string
	version bool
	workers int
}
var errors uint64

func main() {
	loadConfig("config.json")
	parseArgs()

	t := time.Now()
	target := loadTarget()

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
	disconnect()
	if !args.quiet {
		fmt.Printf("%s: %d seconds elapsed\n", args.target, int(time.Since(t).Seconds()))
	}
}

func parseArgs() {
	flag.Uint64Var(&args.errors, "e", config.Errors, "max errors allowed")
	flag.BoolVar(&args.nulls, "n", config.Nulls, "Include nulls in output")
	flag.StringVar(&args.out, "o", "", "Output file or directory")
	flag.IntVar(&args.page, "p", config.Page, "Rows extracted per query")
	flag.BoolVar(&args.quiet, "q", config.Quiet, "Supress informational output")
	flag.Uint64Var(&args.retries, "r", config.Retries, "max consecutive errors")
	flag.BoolVar(&args.version, "v", false, "Print version info about bigboy and exit")
	flag.IntVar(&args.workers, "w", config.Workers, "# of workers")
	flag.Usage = usage
	flag.Parse()

	if args.version {
		stop(version, 0)
	}

	if len(flag.Args()) == 0 {
		usage()
		os.Exit(0)
	}

	args.target = flag.Args()[0]
	args.params = flag.Args()[1:len(flag.Args())]
}

func usage() {
	fmt.Println("usage: bigboy [options] target [params]")
	fmt.Println()
	flag.PrintDefaults()
}
