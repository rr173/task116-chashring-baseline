package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"task116-chashring/internal/api"
	"task116-chashring/internal/selfcheck"
	"task116-chashring/internal/store"
)

func main() {
	var (
		addr        = flag.String("addr", ":8080", "HTTP listen address")
		dbPath      = flag.String("db", "chashring.db", "SQLite database path")
		smoke       = flag.Bool("smoke-test", false, "run self-check and exit")
		migrateOnly = flag.Bool("migrate-only", false, "migrate the store and exit")
	)
	flag.Parse()

	if *smoke {
		if err := selfcheck.Run(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, "smoke-test failed:", err)
			os.Exit(1)
		}
		fmt.Println("smoke-test OK")
		return
	}

	st, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	if *migrateOnly {
		return
	}

	a := api.New(st)
	if err := a.LoadFromStore(context.Background()); err != nil {
		log.Fatalf("load from store: %v", err)
	}

	srv := &http.Server{Addr: *addr, Handler: a.Handler()}
	log.Printf("chashring listening on %s", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}
}
