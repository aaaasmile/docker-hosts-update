package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/aaaasmile/docker-hosts-update/info"
)

func main() {
	log.Println("Welcome to the Hosts updater")
	var dirout = flag.String("dirout", `c:\windows\system32\drivers\etc`, "Test the Hosts output inside a temp file")
	var debug = flag.Bool("debug", false, "Turn on debug information")
	var ver = flag.Bool("ver", false, "Version information")
	flag.Parse()

	if runtime.GOOS != "windows" {
		log.Fatal("Windows only")
	}
	if *ver {
		log.Println("Version 0.2.20201217")
		return
	}

	mmIP, err := info.CollectContainerHostinfo()
	if err != nil {
		log.Fatal(err)
	}
	err = info.UpdateHostsFile(mmIP, *debug, *dirout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Hosts update finished.\nPress any key to finish.")
	fmt.Scanln()

}
