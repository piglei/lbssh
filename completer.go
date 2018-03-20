package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

// MainCompleter is the main completer
type MainCompleter struct {
	hostCompleter HostCompleter
}

func NewMainCompleter(hostCompleter HostCompleter) *MainCompleter {
	return &MainCompleter{
		hostCompleter: hostCompleter,
	}
}

func (cpl *MainCompleter) completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.CurrentLine(), " ")
	if args[0] == ActionGo {
		return cpl.hostCompleter.completer(d)
	}
	return []prompt.Suggest{}
}

// HostCompleter
type HostCompleter struct {
	entris []*HostEntry
}

func (cpl *HostCompleter) completer(d prompt.Document) []prompt.Suggest {
	key := d.GetWordBeforeCursor()
	// Only show suggestions when cursor not in first action word, which means there always will be
	// more than one space before cursor.
	if !strings.Contains(d.TextBeforeCursor(), " ") {
		return []prompt.Suggest{}
	}

	// Allow user to override user, "root@x.com" => "x.com"
	var textBeforekey string
	keySegments := strings.Split(key, "@")
	if len(keySegments) != 1 {
		key = keySegments[len(keySegments)-1]
		textBeforekey = strings.Join(keySegments[:len(keySegments)-1], "@") + "@"
	}

	var suggestions []prompt.Suggest
	for _, hostEntry := range FilterHostsByKeyword(cpl.entris, key) {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        textBeforekey + hostEntry.Name,
			Description: fmt.Sprintf("%s", hostEntry.HostName),
		})
	}
	log.Debugf("%s matches found for key %s", len(suggestions), key)
	return suggestions
}

type MatchedItem struct {
	host   *HostEntry
	weight float64
}

// FilterHostsByKeyword fuzzy query host entries
func FilterHostsByKeyword(entries []*HostEntry, key string) []*HostEntry {
	var matched []*MatchedItem
	for _, hostEntry := range entries {
		nameRank := fuzzy.RankMatchFold(key, hostEntry.Name)
		hostNameRank := fuzzy.RankMatchFold(key, hostEntry.HostName)
		if nameRank == -1 && hostNameRank == -1 {
			continue
		}

		var weight float64
		if nameRank != -1 {
			weight += 1 / float64(nameRank)
		}
		if hostNameRank != -1 {
			weight += 1 / float64(hostNameRank)
		}
		matched = append(matched, &MatchedItem{
			host:   hostEntry,
			weight: weight,
		})
	}
	// Sort results
	sort.SliceStable(matched, func(i, j int) bool {
		if matched[i].weight == matched[j].weight {
			return len(matched[i].host.Name) < len(matched[j].host.Name)
		} else {
			return matched[i].weight > matched[j].weight
		}
	})

	var results []*HostEntry
	for _, item := range matched {
		results = append(results, item.host)
	}
	return results
}
