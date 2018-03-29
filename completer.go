package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"

	"github.com/piglei/lbssh/pkg/storage"
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
	entris  []*HostEntry
	backend storage.HostBackend
}

// HostMatchedItem is the matched item for completer
type HostMatchedItem struct {
	isRecommended bool
	entry         *HostEntry
	profile       *storage.HostProfile
}

func (item *HostMatchedItem) GetDescription() string {
	recommendedFlag := ""
	if item.isRecommended {
		recommendedFlag = "[*] "
	}
	return fmt.Sprintf(
		"%sLast visited: %s | %s",
		recommendedFlag,
		item.profile.GetLastVisitedForDisplay(),
		item.entry.HostName,
	)
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
	var items []*HostMatchedItem
	var mostRecentIndex int
	var maxLastVisited int
	matchedHosts := FilterHostsByKeyword(cpl.entris, key)

	for i, hostEntry := range matchedHosts {
		profile, _ := cpl.backend.GetProfile(hostEntry.Name)
		if profile.LastVisited > maxLastVisited {
			mostRecentIndex = i
			maxLastVisited = profile.LastVisited
		}
		items = append(items, &HostMatchedItem{
			profile: profile,
			entry:   hostEntry,
		})
	}
	// When user input is empty, we should sort all items by last visited time desc
	if key == "" {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].profile.LastVisited > items[j].profile.LastVisited
		})
	} else {
		// Delete the most recent item from sorting
		var mostRecentItem *HostMatchedItem
		if mostRecentIndex != 0 {
			mostRecentItem = items[mostRecentIndex]
			items = append(items[:mostRecentIndex], items[mostRecentIndex+1:]...)

			// Prepend the mostRecentItem back to slice and mark it as RECOMMENDED
			mostRecentItem.isRecommended = true
			itemsCopy := append([]*HostMatchedItem{}, mostRecentItem)
			items = append(itemsCopy, items...)
		}
	}

	for _, item := range items {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        textBeforekey + item.entry.Name,
			Description: item.GetDescription(),
		})
	}
	log.Debugf("%s matches found for key %s", len(suggestions), key)
	return suggestions
}

type MatchedItem struct {
	key           string
	host          *HostEntry
	weight        float64
	numFields     int
	mGroups       []int
	editDistances []int
}

func (m *MatchedItem) SortmGroups() {
	sort.Sort(sort.Reverse(sort.IntSlice(m.mGroups)))
}

// HasHigherPriority compares MatchedItem with another item to determine which one should have higher
// priority in search results.
func (m *MatchedItem) HasHigherPriority(n *MatchedItem) bool {
	if m.numFields != n.numFields {
		return m.numFields > n.numFields
	}
	if len(m.mGroups) != len(n.mGroups) {
		return len(m.mGroups) < len(n.mGroups)
	}

	for i, v := range m.mGroups {
		if v == n.mGroups[i] {
			continue
		}
		return v > n.mGroups[i]
	}

	for i, v := range m.editDistances {
		if v == n.editDistances[i] {
			continue
		}
		return v < n.editDistances[i]
	}
	return true
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
		log.Debugf("Match item found: key=%s numFields=%d groups=%+v", key, matched.numFields, matched.mGroups)
	}

	// Sort results
	sort.SliceStable(matchedItems, func(i, j int) bool {
		m, n := matchedItems[i], matchedItems[j]
		return m.HasHigherPriority(n)
	})

	var results []*HostEntry
	for _, item := range matchedItems {
		results = append(results, item.host)
	}
	return results
}
