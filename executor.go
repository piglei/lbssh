package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"

	"github.com/piglei/lbssh/pkg/storage"
)

type DefaultExecutor struct {
	backend storage.HostBackend
}

func (e *DefaultExecutor) execute(t string) {
	t = strings.TrimSpace(t)
	if t == "" {
		return
	} else if t == "quit" || t == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	args := strings.Split(t, " ")
	if args[0] == ActionGo && len(args) > 1 {
		// Add counter for host
		hostName := args[1]
		err := e.backend.AddNewVisit(hostName)
		if err != nil {
			log.Warnf("Unable to update visit fields for %s: %s", hostName, err.Error())
		}
		sshArgs := strings.Join(args[1:], " ")
		e.StartSSHSession(sshArgs)
		return

	} else {
		log.Infof("Usage: %s HOSTNAME", ActionGo)
		return
	}
	return
}

func (e *DefaultExecutor) StartSSHSession(sshArgs string) {
	sshBin := viper.GetString(SSHBin)
	if sshBin == "" {
		sshBin = SSHBinDefault
	}

	command := sshBin + " " + sshArgs
	log.Infof("SSH into machine using command '%s'", command)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	}
	return
}
