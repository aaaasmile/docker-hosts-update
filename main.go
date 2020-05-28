package main

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	log.Println("Welcome to the Host updater")
	if runtime.GOOS != "windows" {
		log.Fatal("Windows only")
	}

	mm, err := getContainerList()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Recognized container: ", mm)

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

func getContainerList() (map[string]string, error) {
	var cmd string
	var args []string

	cmd = "docker"
	args = []string{"ps", "-q"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Printf("Error on executing docker: %v", err)
		return nil, err
	}

	log.Println("ls -q: ", out)
	list := string(out)
	log.Println("List ", list)

	res := make(map[string]string, 0)
	arr := strings.Split(list, "\n")
	for _, item := range arr {
		res[item] = ""
	}
	return res, nil
}

func getIpDockerContainerIP(contName string) (string, error) {
	// example: docker inspect --format '{{ .NetworkSettings.Networks.nat.IPAddress }}' sql17
	var cmd string
	var args []string

	cmd = "docker"
	args = []string{"inspect", "--format", "'{{ .NetworkSettings.Networks.nat.IPAddress }}'", contName}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Printf("Error on executing docker: %v", err)
		return "", err
	}
	ipstr := strings.Trim(string(out), "'\n")
	log.Println("IP is ", ipstr)

	return ipstr, nil
}
