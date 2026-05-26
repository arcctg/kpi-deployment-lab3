package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	systemdActivation "github.com/coreos/go-systemd/v22/activation"
)

const configPath = "/etc/mywebapp/config.yaml"

func main() {
	migrateOnly := flag.Bool("migrate", false, "run database migration and exit")
	flag.Parse()

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	if err := migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if *migrateOnly {
		log.Println("migration done")
		os.Exit(0)
	}

	app := &App{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handleRoot)
	mux.HandleFunc("/health/alive", app.handleAlive)
	mux.HandleFunc("/health/ready", app.handleReady)
	mux.HandleFunc("/notes", app.handleNotes)
	mux.HandleFunc("/notes/", app.handleNoteByID)

	serve := func(ln net.Listener) {
		srv := &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err := srv.Serve(ln); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}

	listeners, err := systemdActivation.Listeners()
	if err == nil && len(listeners) > 0 {
		log.Printf("starting on systemd socket")
		serve(listeners[0])
		return
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
	log.Printf("starting on %s", addr)
	serve(ln)
}
