package main

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"

	"sshserver/commands"
	"sshserver/utils"
	"sshserver/colors"
	"golang.org/x/term"
	sess "github.com/gliderlabs/ssh"
)

func main() {
	user := "synth"
	address := "192.168.254.152"
	command := "uptime"
	port := "2222"

	key, err := ioutil.ReadFile("/Users/synth/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := kh.New("/Users/synth/.ssh/known_hosts")
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", address+":"+port, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		log.Fatal("unable to create SSH session: ", err)
	}
	defer ss.Close()
	// Creating the buffer which will hold the remotly executed command's output.
	// var stdoutBuf bytes.Buffer
	// ss.Stdout = &stdoutBuf
	ss.Run(command)
	// // Let's print out the result of command.
	// fmt.Println(stdoutBuf.String())

	sess.Handle(func(s sess.Session) {
		utils.Type(s, fmt.Sprintf("%sWelcome to ssh-okra! %s\n\n", colors.Green, colors.Reset))

		term := term.NewTerminal(s, fmt.Sprintf("%s[%s@ssh-okra]%s$ ", colors.Green, s.User(), colors.Reset))

		
				
		for {
			commande, err := term.ReadLine()

			if err != nil {
				CloseServerHandler(s)
				break
			}

			commands.RunCommand(term, commande, s)
			log.Println(fmt.Sprintf("%s: \"%s\"", s.RemoteAddr(), commande))

		}
	})
}

func CloseServerHandler(s sess.Session) {
	

	log.Println(fmt.Sprintf("connection from %s closed.", s.RemoteAddr()))
	// s.Exit(0)
}
