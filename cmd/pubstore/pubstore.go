// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/edrlab/pubstore/pkg/api"
	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/opds"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/edrlab/pubstore/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server context
type Server struct {
	*conf.Config
	*stor.Store
	Router *chi.Mux
}

func main() {

	s := Server{}
	s.Initialize()

	// create an HTTP Server
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", s.Config.Port), Handler: s.Router}

	// server run context
	serverCtx, serverStop := context.WithCancel(context.Background())

	// listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out ... forcing exit.")
			}
		}()

		// trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStop()
	}()

	// run the server
	log.Println("Server starting on port " + strconv.Itoa(s.Config.Port))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// wait for server context to be stopped
	<-serverCtx.Done()

	log.Println("Server halted !!")
}

// Initialize sets the configuration, database and routes
func (s *Server) Initialize() {

	// Initialize the configuration from a config file or/and environment variables
	cfg, err := conf.Init(os.Getenv("PUBSTORE_CONFIG"))
	if err != nil {
		log.Println("Configuration failed: " + err.Error())
		os.Exit(1)
	}
	s.Config = &cfg

	// Initialize the database
	str, err := stor.Init(cfg.DSN)
	if err != nil {
		log.Println("Database setup failed: " + err.Error())
		os.Exit(1)
	}
	s.Store = &str

	// Initialize packages
	_api := api.Init(s.Config, s.Store)
	_view := view.Init(s.Config, s.Store)
	_web := web.Init(s.Config, s.Store, &_view)
	_opds := opds.Init(s.Config, s.Store)

	// Initialize api routes
	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(_api.Router)
	r.Group(_web.Router)
	r.Group(_opds.Router)

	s.Router = r
}
