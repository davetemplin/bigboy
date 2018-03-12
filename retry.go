package main

import (
	"sync/atomic"
	"fmt"
	"time"
)

func retryCheck(err error, retryPtr *uint64) {
	atomic.AddUint64(retryPtr, 1)
	atomic.AddUint64(&errors, 1)

	warn(fmt.Sprintf("error %d of %d", errors, args.errors))
	warn(fmt.Sprintf("%s", err))
	if errors > args.errors {
		stop("error limit exceeded", 2)
	}

	if *retryPtr <= args.retries {
		fmt.Println("pausing for 30 seconds...")	
		time.Sleep(30 * time.Second)
		fmt.Printf("retry %d of %d\n", *retryPtr, args.retries)
	} else {
		stop(fmt.Sprintf("aborted after %d consecutive retries", args.retries), 2)
	}
}