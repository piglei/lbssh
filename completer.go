package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"

	"github.com/piglei/lbssh/pkg/util"
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
	key       string
	host      *HostEntry
	weight    float64
	numFields int
	mGroups   []int
	editDistances []int
}

func (m *MatchedItem) SortmGroups() {
	sort.Sort(sort.Reverse(sort.IntSlice(m.mGroups)))
}

// FilterHostsByKeyword using fuzzy search to find the best matched hostEntried and sort the result
func FilterHostsByKeyword(entries []*HostEntry, key string) []*HostEntry {
	key = strings.ToLower(key)
	var matchedItems []*MatchedItem
	for _, hostEntry := range entries {
		matched := &MatchedItem{
			key:  key,
			host: hostEntry,
		}
		for _, target := range []string{hostEntry.Name, hostEntry.HostName} {
			target = strings.ToLower(target)
			// Use fuzzy match to grep the result first for better performance
			if !fuzzy.Match(key, target) {
				continue
			}
			_, _, mGroups := util.LCSFuzzySearch(key, target)

			log.Debugf("Match field found: key=%s target=%s groups=%+v", key, target, mGroups)
			matched.numFields += 1
			for _, value := range mGroups {
				matched.mGroups = append(matched.mGroups, value)
			}
			// If mGroups is identical, levenshteinDistance will be use for sort
			matched.editDistances = append(matched.editDistances, fuzzy.LevenshteinDistance(key, target))
		}
		// Ignore non-matched found items
		if matched.numFields == 0 {
			continue
		}

		matched.SortmGroups()
		matchedItems = append(matchedItems, matched)
		log.Debugf("Match item found: key=%s groups=%+v", key, matched.mGroups)
	}

	// Sort results
	sort.SliceStable(matchedItems, func(i, j int) bool {
		i1, i2 := matchedItems[i], matchedItems[j]
		if i1.numFields > i2.numFields {
			return true
		}
		if len(i1.mGroups) != len(i2.mGroups) {
			return len(i1.mGroups) < len(i2.mGroups)
		}

		for i, v := range i1.mGroups {
			if v == i2.mGroups[i] {
				continue
			}
			return v > i2.mGroups[i]
		}

		for i, v := range i1.editDistances {
			if v == i2.editDistances[i] {
				continue
			}
			return v < i2.editDistances[i]
		}
		return true
	})

	var results []*HostEntry
	for _, item := range matchedItems {
		results = append(results, item.host)
	}
	return results
}
