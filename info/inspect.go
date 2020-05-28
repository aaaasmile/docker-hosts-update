package info

import (
	"log"
	"os/exec"
	"strings"
)

func CollectContainerHostinfo() (map[string]string, error) {
	mmName, err := getContainerList()
	if err != nil {
		return nil, err
	}
	log.Println("Recognized container: ", mmName)

	mmIP := make(map[string]string)
	for _, name := range mmName {
		ip, err := getIpDockerContainerIP(name)
		if err != nil {
			return nil, err
		}
		mmIP[name] = ip
	}
	return mmIP, nil
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
