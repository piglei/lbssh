package main

import (
	"regexp"
	"strings"
)

type HostEntry struct {
	Name       string
	HostName   string
	Annotation string
}

type SSHConfigFileParser struct {
	Content string

	hostEntries      []*HostEntry
	currentHostEntry *HostEntry
}

func NewSSHConfigFileParser(content string) *SSHConfigFileParser {
	return &SSHConfigFileParser{
		Content: content,
	}
}

// Parse current file content, after calling Parse, use .Result() to get all HostEntries.
func (parser *SSHConfigFileParser) Parse() {
	var ReBlank = regexp.MustCompile(`\s+`)
	var ReIPAddress = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

	for _, line := range strings.Split(parser.Content, "\n") {
		items := ReBlank.Split(strings.TrimSpace(line), 2)
		if len(items) < 2 {
			continue
		}
		key, value := items[0], items[1]
		// SSH config file keyword is case-insensitive
		switch strings.ToLower(key) {
		case "hostname":
			if parser.currentHostEntry != nil {
				parser.currentHostEntry.HostName = value
			}
		case "host":
			// End of host section
			if parser.currentHostEntry != nil {
				parser.hostEntries = append(parser.hostEntries, parser.currentHostEntry)
				parser.currentHostEntry = nil
			}
			if strings.Index(value, "*") != -1 {
				break
			}
			parser.currentHostEntry = &HostEntry{Name: value}
		case "proxycommand":
			// If Host was connected via a proxy command, try to get the IP address as HostName
			IPInCommand := ReIPAddress.Find([]byte(value))
			if IPInCommand != nil {
				if parser.currentHostEntry != nil {
					parser.currentHostEntry.HostName = string(IPInCommand)
				}
			}
		}
	}
	// Append the last HostEntry
	if parser.currentHostEntry != nil {
		parser.hostEntries = append(parser.hostEntries, parser.currentHostEntry)
	}
}

func (parser *SSHConfigFileParser) Result() []*HostEntry {
	return parser.hostEntries
}
