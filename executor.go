package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

func executor(t string) {
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
		sshBin := viper.GetString(SSHBin)
		if sshBin == "" {
			sshBin = SSHBinDefault
		}

		command := sshBin + " " + strings.Join(args[1:], " ")
		log.Infof("SSH into machine using command '%s'", command)
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
		}
		return
	} else {
		log.Infof("Usage: %s HOSTNAME", ActionGo)
		return
	}
}
