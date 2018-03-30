package storage

import (
	"fmt"
	"github.com/asdine/storm"
	"time"
)

type HostProfile struct {
	ID          int    `storm:"id,increment"`
	Name        string `storm:"unique"`
	Visited     int
	LastVisited int
}

func (p *HostProfile) GetLastVisitedForDisplay() string {
	if p.LastVisited == 0 {
		return "-"
	}
	return RelativeTimeDisplay(p.LastVisited)
}

func RelativeTimeDisplay(ts int) string {
	tsDelta := int(time.Now().Unix()) - ts
	if tsDelta < 60 {
		return fmt.Sprintf("%ds", tsDelta)
	} else if tsDelta < 3600 {
		return fmt.Sprintf("%dm", tsDelta/60)
	} else if tsDelta < 3600*24 {
		return fmt.Sprintf("%dh", tsDelta/3600)
	} else if tsDelta < 3600*24*7 {
		return fmt.Sprintf("%dd", tsDelta/(3600*24))
	} else {
		return "7d+"
	}
	return "unknown"
}

// HostBackend defines the interface for hostEntry actions
type HostBackend interface {
	Open() error
	Close() error
	CreateProfile(hostname string) (*HostProfile, error)
	GetProfile(hostname string) (*HostProfile, error)
	AddNewVisit(hostname string) error
	DeleteHost(hostname string) error
}

type HostBackendStorm struct {
	dbPath string
	Database *storm.DB
}

func NewHostBackend(path string) (*HostBackendStorm, error) {
	return &HostBackendStorm{
		dbPath: path,
	}, nil
}
func (backend *HostBackendStorm) Open() error {
	db, err := storm.Open(backend.dbPath)
	if err != nil {
		return err
	}
	backend.Database = db
	return nil
}

func (backend *HostBackendStorm) Close() error {
	if backend.Database != nil {
		return backend.Database.Close()
	}
	return nil
}

func (backend *HostBackendStorm) CreateProfile(hostname string) (*HostProfile, error) {
	profile := HostProfile{Name: hostname}
	if err := backend.Database.Save(&profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (backend *HostBackendStorm) GetProfile(hostname string) (*HostProfile, error) {
	var profile HostProfile
	err := backend.Database.One("Name", hostname, &profile)
	return &profile, err
}

func (backend *HostBackendStorm) AddNewVisit(hostname string) error {
	profile, err := backend.GetProfile(hostname)
	if err != nil {
		return err
	}
	profile.Visited++
	profile.LastVisited = int(time.Now().Unix())
	return backend.Database.Update(profile)
}

func (backend *HostBackendStorm) DeleteHost(hostname string) error {
	profile, err := backend.GetProfile(hostname)
	if err != nil {
		return err
	}
	return backend.Database.DeleteStruct(profile)
}
