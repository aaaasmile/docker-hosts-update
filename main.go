package main

import (
	"log"
	"runtime"
)

func main() {
	log.Println("Welcome to the Host updater")
	if runtime.GOOS != "windows" {
		log.Fatal("Windows only")
	}

	mmIP, err := info.CollectContainerHostinfo()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Now check if the Hosts file need to be updated with ", mmIP)
}
