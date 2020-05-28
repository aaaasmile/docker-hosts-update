package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	log.Println("Welcome to the Host updater")

	ip1, err := getIpDockerContainerIP("a266530d48f4")
	if err != nil {
		log.Fatal(err)
	}
	ip2, err := getIpDockerContainerIP("21c9649b5768")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Now check if the Hosts file need to be updated with ", ip1, ip2)
}

func getIpDockerContainerIP(contName string) (string, error) {
	// example: docker inspect --format '{{ .NetworkSettings.Networks.nat.IPAddress }}' sql17
	var cmd string
	var args []string

	if runtime.GOOS == "windows" {
		cmd = "docker"
		args = []string{"inspect", "--format", "'{{ .NetworkSettings.Networks.nat.IPAddress }}'", contName}
	} else {
		log.Println("OS not recognized")
		return "", fmt.Errorf("OS not supported %s", runtime.GOOS)
	}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Printf("Error on executing docker: %v", err)
		return "", err
	}
	ipstr := strings.Trim(string(out), "'\n")
	log.Println("IP is ", ipstr)

	return ipstr, nil
}
