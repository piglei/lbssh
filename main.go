package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"text/template"

	"bytes"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/piglei/lbssh/pkg/storage"
	"github.com/piglei/lbssh/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ActionGo = "go"

	SSHBin             = "SSH_BIN"
	SSHBinDefault      = "/usr/bin/ssh"
	WelcomeMessageTmpl = `
  _ _          _
 | | |__ _____| |_
 | | '_ (_-<_-< ' \
 |_|_.__/__/__/_||_|

Welcome to lbssh! (version: {{.GetVersion}})
Please input "go HOSTNAME" to start!`
)

func main() {
	currentUser, _ := user.Current()
	pflag.String("ssh-config-file", path.Join(currentUser.HomeDir, "/.ssh/config"), "ssh config file location")
	pflag.String("storage-db-file", path.Join(currentUser.HomeDir, "/.lbssh.db"), "db file location")
	pflag.String("ssh-bin", SSHBinDefault, "ssh binary path")
	pflag.String("log-level", "INFO", "log level")
	pflag.Bool("version", false, "display version info")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	// Allow set configs from environment variables
	viper.SetEnvPrefix("LBSSH")
	viper.BindEnv(SSHBin)

	if viper.GetBool("version") {
		currentVersion := version.Get()
		fmt.Printf("lbssh version: {%s}\n", currentVersion.ForDisplay())
		os.Exit(0)
	}

	// Setting log-level
	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatalf("Unable to set log level: %s", err.Error())
	}
	log.SetLevel(level)
	log.Debugf("Log level was set to %s", log.GetLevel())

	fnameConfig := viper.GetString("ssh-config-file")
	sshConfigContent, err := ioutil.ReadFile(fnameConfig)
	if err != nil {
		log.Fatalf("unable to open ssh config file: %s\n", err)
		os.Exit(1)
	}

	parser := NewSSHConfigFileParser(string(sshConfigContent))
	parser.Parse()

	// Update hostProfiles in storage
	storageDBFile := viper.GetString("storage-db-file")
	if storageDBFile == "" {
		log.Fatalf("storage db file path can't be empty.\n")
		os.Exit(1)
	}
	backend, err := storage.NewHostBackend(storageDBFile)
	if err != nil {
		log.Fatalf("Unable to open db file: %s\n", err.Error())
		os.Exit(2)
	}

	err = backend.Open()
	if err != nil {
		log.Fatalf("Unable to open database: %s\n", err.Error())
		os.Exit(2)
	}
	// Initialize all hostProfiles
	for _, host := range parser.Result() {
		backend.CreateProfile(host.Name)
	}
	backend.Close()

	hostCompleter := HostCompleter{
		entris:  parser.Result(),
		backend: backend,
	}
	mainCompleter := NewMainCompleter(hostCompleter)

	fmt.Println(GetWelcomeMessage())
	e := DefaultExecutor{backend: backend}
	p := prompt.New(
		e.execute,
		mainCompleter.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionSwitchKeyBindMode(prompt.EmacsKeyBind),
		prompt.OptionMaxSuggestion(6),
	)
	p.Run()
}

// GetWelcomeMessage returns the welcome message when user logged in
func GetWelcomeMessage() string {
	tmpl, _ := template.New("welcome_message").Parse(WelcomeMessageTmpl)
	result := bytes.Buffer{}
	tmpl.Execute(&result, version.Get())
	return result.String()
}
