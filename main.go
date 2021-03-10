package main

import (
	"bytes"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

// Args - Command Line Arguments
type Args struct {
	config  string
	errors  uint64
	nulls   bool
	out     string
	page    int
	params  []string
	quiet   bool
	retries uint64
	target  string
	version bool
	workers int
}

var args Args
var errors uint64

func main() {
	config = loadConfig(defaultConfig)
	args, _ := parseArgs(os.Args[0], os.Args[1:], config)

	if args.version {
		stop(version, 0)
	}

	if args.target == "" {
		usage()
		os.Exit(0)
	}

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

func parseArgs(progname string, input []string, c Config) (args *Args, output string) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var a Args
	flags.StringVar(&a.config, "c", defaultConfig, "Bigboy conifg file path")
	flags.Uint64Var(&a.errors, "e", c.Errors, "max errors allowed")
	flags.BoolVar(&a.nulls, "n", c.Nulls, "Include nulls in output")
	flags.StringVar(&a.out, "o", "", "Output file or directory")
	flags.IntVar(&a.page, "p", c.Page, "Rows extracted per query")
	flags.BoolVar(&a.quiet, "q", c.Quiet, "Supress informational output")
	flags.Uint64Var(&a.retries, "r", c.Retries, "max consecutive errors")
	flags.BoolVar(&a.version, "v", false, "Print version info about bigboy and exit")
	flags.IntVar(&a.workers, "w", c.Workers, "# of workers")
	flags.Usage = usage

	err := flags.Parse(input)
	if err != nil {
		fmt.Println("Error reading arguments", err)
		os.Exit(0)
	}

	allArgs := flags.Args()
	if len(allArgs) == 0 {
		return &a, buf.String()
	}

	a.target = allArgs[0]
	a.params = allArgs[1:]
	return &a, buf.String()
}

func usage() {
	fmt.Println("usage: bigboy [options] target [params]")
	fmt.Println()
	flag.PrintDefaults()
}
