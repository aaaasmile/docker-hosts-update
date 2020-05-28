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

	mmName, err := getContainerList()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Recognized container: ", mmName)

	mmIp := make(map[string]string)
	for _, name := range mmName {
		ip, err := getIpDockerContainerIP(name)
		if err != nil {
			log.Fatal(err)
		}
		mmIp[name] = ip
	}

	log.Println("Now check if the Hosts file need to be updated with ", mmIp)
}

func getContainerList() ([]string, error) {
	var cmd string
	var args []string

	cmd = "docker"
	args = []string{"ps", "-q"}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Printf("Error on executing docker: %v", err)
		return nil, err
	}

	//fmt.Println("*** ls -q: ", out)
	list := string(out)
	log.Println("List ", list)

	res := make([]string, 0)
	arr := strings.Split(list, "\n")
	for _, contHash := range arr {
		if len(contHash) > 0 {
			args = []string{"inspect", "--format", "'{{ .Name }}'", contHash}
			log.Println("Inspect container ", args)
			out, err := exec.Command(cmd, args...).Output()
			if err != nil {
				log.Printf("Error on executing docker: %v", err)
				return nil, err
			}
			res = append(res, strings.Trim(string(out), "/'\n"))
		}
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
