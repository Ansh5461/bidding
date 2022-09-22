package api

import (
	"errors"
	"net/http"
	"time"

	pkg "bidding/pkg"
	"bidding/store"

	"github.com/gorilla/context"
)

// Store holds new conn
var Store *store.Conn

// ServiceInfo stores basic service information
type ServiceInfo struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Uptime  time.Time `json:"uptime"`
	Epoch   int64     `json:"epoch"`
}

// ServiceName holds the service which connected to
var ServiceName = ""
var serviceInfo *ServiceInfo

// InitService sets the service name
func InitService(name, version string) {
	ServiceName = name
	serviceInfo = &ServiceInfo{
		Name:    name,
		Version: version,
		Uptime:  time.Now(),
		Epoch:   time.Now().Unix(),
	}

	Store = store.NewStore()
	// bidder.ShareConn()
}

// API Handler's ---------------------------------------------------------------

// Handler custom api handler help us to handle all the errors in one place
type Handler func(w http.ResponseWriter, r *http.Request) *string

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	// clear gorilla context
	defer context.Clear(r)
	if err != nil {
		// APP Level Error
		pkg.Log.Infof("ServiceName: %s, StatusCode: %d, Error: %s\n ",
			ServiceName, http.StatusInternalServerError, errors.New("got error while serving http"))
		pkg.Fail(w, errors.New("got error while serving http"))
	}
}

// Basic Handler func ---------------------------------------------------------------

// IndexHandeler common index handler for all the service
func IndexHandeler(w http.ResponseWriter, r *http.Request) {
	pkg.OK(w, map[string]string{
		"name":    serviceInfo.Name,
		"version": serviceInfo.Version,
	})
}

// HealthHandeler return basic service info
func HealthHandeler(w http.ResponseWriter, r *http.Request) {
	pkg.OK(w, serviceInfo)
}
