package main_test

import (
	. "github.com/piglei/lbssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterHostsByKeyword", func() {
	It("test proper order", func() {
		hosts := []*HostEntry{
			{Name: "capggp"},
			{Name: "apple"},
			{Name: "3223abbp3223pppa"},
			{Name: "capgp"},
			{Name: "non-sense"},
		}

		result := FilterHostsByKeyword(hosts, "app")
		Expect(len(result)).To(Equal(4))
		Expect(result).To(Equal(
			[]*HostEntry{
				{Name: "apple"},
				{Name: "capgp"},
				{Name: "capggp"},
				{Name: "3223abbp3223pppa"},
			}))

		result = FilterHostsByKeyword(hosts, "non-existed-key")
		Expect(len(result)).To(Equal(0))
	})

	It("test proper order with Hostname", func() {
		hosts := []*HostEntry{
			{Name: "proxy3-80", HostName: "201.222.222.80"},
			{Name: "q.dev1-180.uni", HostName: "200.222.222.180"},
		}

		result := FilterHostsByKeyword(hosts, "180")
		Expect(result).To(Equal(
			[]*HostEntry{
				{Name: "q.dev1-180.uni", HostName: "200.222.222.180"},
				{Name: "proxy3-80", HostName: "201.222.222.80"},
			}))
	})
})
