package main

import (
	"golang.org/x/crypto/ssh"
	// "golang.org/x/crypto/ssh/knownhosts"
	"io"
	"log"
	"os"
	"io/ioutil"
	// "fmt"
)

func main() {

	// Create client config
	// var hostt = ""
	// var users = ""

	// fmt.Print("host: ")
	// fmt.Scanf("%s\n", &hostt)
	// fmt.Print("user: ")
	// fmt.Scanf("%s\n", &users)

	key, err := ioutil.ReadFile("/Users/synth/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	signer, errs := ssh.ParsePrivateKey(key)
	if errs != nil {
		log.Fatalf("unable to parse private key: %v", errs)
	}
	// var hostkeyCallback ssh.HostKeyCallback
	// hostkeyCallback, err = knownhosts.New("/Users/synth/.ssh/known_hosts")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	config := &ssh.ClientConfig{
		User: "synth",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "192.168.254.152:2222", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("vt220", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run("ls -laG /")
}
