package storage_test

import (
	. "github.com/piglei/lbssh/pkg/storage"

	"crypto/rand"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/big"
	"time"
)

var _ = Describe("Test HostBackendStorm", func() {
	It("integrated test", func() {
		randInt, _ := rand.Int(rand.Reader, big.NewInt(1000))
		hostName := fmt.Sprintf("x-%d.com", randInt)

		backend, _ := NewHostBackend("/tmp/lbssh_test.db")
		backend.CreateProfile(hostName)

		profile, _ := backend.GetProfile(hostName)
		Expect(profile.Visited).To(Equal(0))

		err := backend.AddNewVisit(hostName)
		Expect(err).To(BeNil())

		profile, _ = backend.GetProfile(hostName)
		Expect(profile.Visited).To(Equal(1))
		Expect(profile.LastVisited > 0).To(BeTrue())

		fmt.Println(profile.GetLastVisitedForDisplay())

		err = backend.DeleteHost(hostName)
		Expect(err).To(BeNil())
		_, err = backend.GetProfile(hostName)
		Expect(err).NotTo(BeNil())

	})
})

var _ = Describe("Test RelativeTimeDisplay", func() {
	It("normal tests", func() {
		now := int(time.Now().Unix())
		Expect(RelativeTimeDisplay(now - 5)).To(Equal("5s"))
		Expect(RelativeTimeDisplay(now - 61)).To(Equal("1m"))
		Expect(RelativeTimeDisplay(now - 1210)).To(Equal("20m"))
		Expect(RelativeTimeDisplay(now - 4000)).To(Equal("1h"))
		Expect(RelativeTimeDisplay(now - 3600*21)).To(Equal("21h"))
		Expect(RelativeTimeDisplay(now - 3600*24)).To(Equal("1d"))
		Expect(RelativeTimeDisplay(now - 3600*24*6)).To(Equal("6d"))
		Expect(RelativeTimeDisplay(now - 3600*24*30)).To(Equal("7d+"))
	})
})
