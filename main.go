package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/aaaasmile/docker-hosts-update/info"
)

func main() {
	log.Println("Welcome to the Hosts updater")
	var test = flag.Bool("test", false, "Test the Hosts output inside a temp file")
	var debug = flag.Bool("debug", false, "Turn on debug information")
	flag.Parse()

	if runtime.GOOS != "windows" {
		log.Fatal("Windows only")
	}

	mmIP, err := info.CollectContainerHostinfo()
	if err != nil {
		log.Fatal(err)
	}
	err = info.UpdateHostsFile(mmIP, *debug, *test)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Hosts update finished")
}
