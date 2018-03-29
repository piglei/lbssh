package main_test

import (
	. "github.com/piglei/lbssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

var _ = Describe("FilterHostsByKeyword", func() {
	It("a normal match for multiple cases", func() {
		hosts := []*HostEntry{
			{Name: "capggp"},
			{Name: "apple-2"},
			{Name: "3223abbp3223pppa"},
			{Name: "capgp"},
			{Name: "non-sense"},
		}

		result := FilterHostsByKeyword(hosts, "app")
		Expect(len(result)).To(Equal(4))
		Expect(result).To(Equal(
			[]*HostEntry{
				{Name: "apple-2"},
				{Name: "capgp"},
				{Name: "capggp"},
				{Name: "3223abbp3223pppa"},
			}))

		result = FilterHostsByKeyword(hosts, "non-existed-key")
		Expect(len(result)).To(Equal(0))
	})
	It("key appears in both name & hostname is better", func() {
		hosts := []*HostEntry{
			{Name: "tokyo-180", HostName: "ip-tokyo-180"},
			{Name: "tokyo-jp-191", HostName: "111.111.111.111"},
		}

		result := FilterHostsByKeyword(hosts, "tok")
		Expect(result).To(Equal(
			[]*HostEntry{
				{Name: "tokyo-180", HostName: "ip-tokyo-180"},
				{Name: "tokyo-jp-191", HostName: "111.111.111.111"},
			}))
	})
	It("full segment match is better", func() {
		hosts := []*HostEntry{
			{Name: "cb--", HostName: ""},
			{Name: "ceb", HostName: ""},
		}
		result := FilterHostsByKeyword(hosts, "cb")
		Expect(result).To(Equal(
			[]*HostEntry{
				{Name: "cb--", HostName: ""},
				{Name: "ceb", HostName: ""},
			}))
	})
	It("mGroups is the same, use edit distance", func() {
		hosts := []*HostEntry{
			{Name: "apple-tree-fool", HostName: ""},
			{Name: "apple-tree-in-forest", HostName: ""},
		}
		result := FilterHostsByKeyword(hosts, "appfo")
		Expect(result).To(Equal(
			[]*HostEntry{
				{Name: "apple-tree-fool", HostName: ""},
				{Name: "apple-tree-in-forest", HostName: ""},
			}))
	})
})
