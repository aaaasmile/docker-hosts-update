package main

import (
	"log"
	"runtime"

	"github.com/aaaasmile/docker-hosts-update/info"
)

func main() {
	log.Println("Welcome to the Hosts updater")
	if runtime.GOOS != "windows" {
		log.Fatal("Windows only")
	}

	mmIP, err := info.CollectContainerHostinfo()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Now check if the Hosts file needs to be updated with ", mmIP)
}
