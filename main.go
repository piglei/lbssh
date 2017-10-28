package main

import (
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const (
	ActionGo = "go"

	SSH_BIN         = "SSH_BIN"
	SSH_BIN_DEFAULT = "/usr/bin/ssh"
)

func main() {
	log.SetLevel(log.InfoLevel)

	viper.SetEnvPrefix("LBSSH")
	viper.BindEnv(SSH_BIN)

	currentUser, _ := user.Current()
	// TODO: Support more config file location
	SSHConfigFile := path.Join(currentUser.HomeDir, "/.ssh/config")
	fnameConfig := SSHConfigFile

	sshConfigContent, err := ioutil.ReadFile(fnameConfig)
	if err != nil {
		log.Fatalf("unable to open ssh config file: %s\n", err)
		os.Exit(1)
	}

	parser := NewSSHConfigFileParser(string(sshConfigContent))
	parser.Parse()

	hostCompleter := HostCompleter{
		entris: parser.Result(),
	}
	mainCompleter := NewMainCompleter(hostCompleter)

	p := prompt.New(
		executor,
		mainCompleter.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionSwitchKeyBindMode(prompt.EmacsKeyBind),
		prompt.OptionMaxSuggestion(6),
	)
	p.Run()
}
