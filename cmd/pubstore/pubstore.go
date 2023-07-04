package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edrlab/pubstore/pkg/api"
	"github.com/edrlab/pubstore/pkg/config"
	"github.com/edrlab/pubstore/pkg/opds"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/edrlab/pubstore/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("no .env file found")
	}
	config.Init()

	_stor := stor.Init(config.DSN)
	_api := api.Init(_stor)
	_view := view.Init(_stor)
	_web := web.Init(_stor, _view)
	_opds := opds.Init(_stor)

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(_api.Rooter)
	r.Group(_web.Rooter)
	r.Group(_opds.Router)

	// The HTTP Server
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", config.PORT), Handler: r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Println("Server started on port " + fmt.Sprintf("%d", config.PORT))
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	_stor.Stop()

	log.Println("Server halted !!")
}
