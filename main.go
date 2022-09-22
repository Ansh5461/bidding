package main

import (
	"bidding/api"
	v1 "bidding/api/v1"
	"bidding/config"
	"bidding/pkg"
	"bidding/store"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

type ServiceInfo struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Uptime  time.Time `json:"uptime"`
	Epoch   int64     `json:"epoch"`
}

type Bidder struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Host  string        `json:"host"`
	Delay time.Duration `json:"delay"`
}

var Store *store.Conn

type Bidders interface {
	Add(bidder Bidder)
	List() []Bidder
	Count() int
}

var (
	name        = "bidder"
	version     = "1.0.0"
	ServiceName = ""
	serviceInfo *ServiceInfo
)

func main() {
	ServiceName = name
	serviceInfo = &ServiceInfo{
		Name:    name,
		Version: version,
		Uptime:  time.Now(),
		Epoch:   time.Now().Unix(),
	}

	Store = store.NewStore()

	pkg.Setup(config.Env)

	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{
			"Origin", "Authorization", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Header", "Accept",
			"Content-Type", "X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin", "Origin",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(c.Handler)

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Get("/", api.IndexHandeler)
	r.Get("/top", api.HealthHandeler)
	r.Route("/v1", v1.Init)

}
