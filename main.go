package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

// Args ...
type Args struct {
	errors  uint64
	nulls   bool
	out     string
	page    int
	params  []string
	quiet   bool
	retries uint64
	target  string
	workers int
}

// Flags - Command Line Arguments
type Flags struct {
	Args
	config  string
	version bool
}

var args Args
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

func setFlags(flags *flag.FlagSet, c *Config, f *Flags) {
	flags.StringVar(&f.config, "c", defaultConfig, "Bigboy conifg file path")
	flags.Uint64Var(&f.errors, "e", c.Errors, "max errors allowed")
	flags.BoolVar(&f.nulls, "n", c.Nulls, "Include nulls in output")
	flags.StringVar(&f.out, "o", "", "Output file or directory")
	flags.IntVar(&f.page, "p", c.Page, "Rows extracted per query")
	flags.BoolVar(&f.quiet, "q", c.Quiet, "Supress informational output")
	flags.Uint64Var(&f.retries, "r", c.Retries, "max consecutive errors")
	flags.BoolVar(&f.version, "v", false, "Print version info about bigboy and exit")
	flags.IntVar(&f.workers, "w", c.Workers, "# of workers")
}

func getFlagsUsage(f *flag.FlagSet) func() {
	return func() {
		fmt.Println("usage: bigboy [options] target [params]")
		fmt.Println()
		f.PrintDefaults()
	}
}

func parseFlags(progname string, input []string, c *Config) (*Flags, func()) {
	var flags *flag.FlagSet = flag.NewFlagSet(progname, flag.ContinueOnError)

	var f Flags
	setFlags(flags, c, &f)
	var flagsUsage func() = getFlagsUsage(flags)
	flags.Usage = flagsUsage

	err := flags.Parse(input)
	if err != nil {
		fmt.Println("Error reading arguments", err)
		os.Exit(0)
	}

	return &f, flagsUsage
}

func parseArgs(progname string, input []string, c *Config) *Args {
	var flags *flag.FlagSet = flag.NewFlagSet(progname, flag.ContinueOnError)

	var f Flags
	setFlags(flags, c, &f)

	err := flags.Parse(input)
	if err != nil {
		fmt.Println("Error reading arguments", err)
		os.Exit(0)
	}

	allArgs := flags.Args()
	if len(allArgs) == 0 {
		return &f.Args
	}

	f.target = allArgs[0]
	f.params = allArgs[1:]
	return &f.Args
}
