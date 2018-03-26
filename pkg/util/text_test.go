package util_test

import (
	. "github.com/piglei/lbssh/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LCSFuzzySearch", func() {
	It("search with exactly match", func() {
		matchedLen, matchedString, matchedGroups := LCSFuzzySearch("apple", "apple")
		Expect(matchedLen).To(Equal(5))
		Expect(matchedString).To(Equal("apple"))
		Expect(matchedGroups).To(Equal([]int{5}))
	})
	It("search with exactly match", func() {
		matchedLen, matchedString, matchedGroups := LCSFuzzySearch("2gffm3wxg", "q.ffm-host-wxstag.1")
		Expect(matchedLen).To(Equal(6))
		Expect(matchedString).To(Equal("ffmwxg"))
		Expect(matchedGroups).To(Equal([]int{3, 2, 1}))
	})
	It("search with no match", func() {
		matchedLen, matchedString, matchedGroups := LCSFuzzySearch("apple", "zoo")
		Expect(matchedLen).To(Equal(0))
		Expect(matchedString).To(Equal(""))
		Expect(matchedGroups).To(Equal([]int{}))
	})
	It("search with source longer", func() {
		matchedLen, matchedString, matchedGroups := LCSFuzzySearch("apple", "ppe")
		Expect(matchedLen).To(Equal(3))
		Expect(matchedString).To(Equal("ppe"))
		Expect(matchedGroups).To(Equal([]int{3}))
	})
})
