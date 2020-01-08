package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	workers := flag.Uint("c", 10, "number of concurrent requests to make")
	ips := flag.String("f", "", "file containing list of ips, defaults to stdin if omitted")
	verbose := flag.Bool("v", false, "log errors to stdout")
	flag.Parse()

	var input io.ReadCloser
	if *ips == "" {
		input = os.Stdin
	} else {
		f, err := os.Open(*ips)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open input file: %v\n", err)
			os.Exit(1)
		}
		input = f
	}
	defer input.Close()

	work := make(chan string)
	go makeWork(input, work, *verbose)

	wg := &sync.WaitGroup{}
	for i := 0; i < int(*workers); i++ {
		wg.Add(1)
		go doLookups(wg, work, *verbose)
	}
	wg.Wait()
}

func makeWork(input io.Reader, work chan string, verbose bool) {
	s := bufio.NewScanner(input)
	for s.Scan() {
		work <- s.Text()
	}
	if s.Err() != nil && verbose {
		fmt.Fprintf(os.Stderr, "error while scanning input: %v\n", s.Err())
	}
	close(work)
}

func doLookups(wg *sync.WaitGroup, work chan string, verbose bool) {
	for ip := range work {
		hosts, err := net.LookupAddr(ip)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "error resolving %s: %v\n", ip, err)
			}
			continue
		}
		for _, v := range hosts {
			fmt.Println(strings.TrimSuffix(v, "."))
		}
	}
	wg.Done()
}
