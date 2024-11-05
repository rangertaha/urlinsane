package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/d3mondev/resolvermt"
)

func main() {
	var resolverFile, domainFile, cpuprofileFile, memprofileFile string
	var retries, qps, concurrency int
	flag.StringVar(&resolverFile, "resolvers", "", "filename containing resolvers")
	flag.StringVar(&domainFile, "domains", "", "filename containing domains to resolve")
	flag.StringVar(&cpuprofileFile, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&memprofileFile, "memprofile", "", "write memory profile to file")
	flag.IntVar(&retries, "retry", 3, "number of times to retry a DNS query")
	flag.IntVar(&qps, "qps", 10, "maximum number of queries per second for each resolver")
	flag.IntVar(&concurrency, "concurrency", 1000, "number of concurrent DNS requests")
	flag.Parse()

	if resolverFile == "" || domainFile == "" {
		usage()
		os.Exit(1)
	}

	if cpuprofileFile != "" {
		file, err := os.Create(cpuprofileFile)

		if err != nil {
			fmt.Printf("unable to write cpu profile file %s\n\n", cpuprofileFile)
			os.Exit(1)
		}

		if err := pprof.StartCPUProfile(file); err != nil {
			fmt.Printf("unable to start cpu profiling\n\n")
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	content, err := ioutil.ReadFile(resolverFile)

	if err != nil {
		fmt.Printf("unable to open file %s\n\n", resolverFile)
		usage()
		os.Exit(1)
	}

	resolvers := strings.Split(string(content), "\n")

	content, err = ioutil.ReadFile(domainFile)

	if err != nil {
		fmt.Printf("unable to open file %s\n\n", domainFile)
		usage()
		os.Exit(1)
	}

	domains := strings.Split(string(content), "\n")

	client := resolvermt.New(resolvers, retries, qps, concurrency)
	defer client.Close()

	records := client.Resolve(domains, resolvermt.TypeA)

	for _, record := range records {
		fmt.Printf("%s %s %s\n", record.Question, record.Type, record.Answer)
	}

	if memprofileFile != "" {
		file, err := os.Create(memprofileFile)

		if err != nil {
			fmt.Printf("unable to write memory profile file %s\n\n", memprofileFile)
			os.Exit(1)
		}

		defer file.Close()

		if err := pprof.WriteHeapProfile(file); err != nil {
			fmt.Printf("unable to write memory profile\n\n")
			os.Exit(1)
		}
	}
}

func usage() {
	fmt.Printf("Usage: %s --resolvers <file> --domains <file> [flags]\n", os.Args[0])
	flag.PrintDefaults()
}
