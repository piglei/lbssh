package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"

	"github.com/piglei/lbssh/pkg/util"
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
	host      *HostEntry
	weight    float64
	numFields int
	mGroups   []int
}

func (m *MatchedItem) SortmGroups() {
	sort.Sort(sort.Reverse(sort.IntSlice(m.mGroups)))
}

// FilterHostsByKeyword using fuzzy search to find the best matched hostEntried and sort the result
func FilterHostsByKeyword(entries []*HostEntry, key string) []*HostEntry {
	// Because we are using Levenshtein distance(aka. edit distance) algorithm to find the best matched result,
	// in order to perform an apples to apples comparison, we can:
	//
	// - Remove the leading and trailing characters, then count for edit distance
	var maxLengthName, maxLengthHostname int
	for _, hostEntry := range entries {
		if len(hostEntry.Name) > maxLengthName {
			maxLengthName = len(hostEntry.Name)
		}
		if len(hostEntry.HostName) > maxLengthHostname {
			maxLengthHostname = len(hostEntry.HostName)
		}
	}
	log.Debugf(
		"maxLengthName=%d, maxLengthHostname=%d, will extend all cadidates to this length",
		maxLengthName,
		maxLengthHostname,
	)

	var matchedItems []*MatchedItem
	for _, hostEntry := range entries {
		matched := &MatchedItem{host: hostEntry}
		for _, target := range []string{hostEntry.Name, hostEntry.HostName} {
			mLength, _, mGroups := util.LCSFuzzySearch(key, target)
			if mLength != len(key) {
				continue
			}

			log.Debugf("Match field found: key=%s target=%s groups=%+v", key, target, mGroups)
			matched.numFields += 1
			for _, value := range mGroups {
				matched.mGroups = append(matched.mGroups, value)
			}
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
		return true
	})

	var results []*HostEntry
	for _, item := range matchedItems {
		results = append(results, item.host)
	}
	return results
}
