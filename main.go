package main

import (
	"flag"
	"log"

	"github.com/crowdworks/cyqldog/cyqldog"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	// Version information embedded at build time.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log.Printf("main: starting cyqldog (version: %s, commit: %s, date: %s)", version, commit, date)

	// Parse the argument's flag.
	var configPath string
	flag.StringVar(&configPath, "C", "./cyqldog.yml", "path to config file")
	flag.Parse()

	m := cyqldog.NewMonitor(configPath)
	if err := m.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("main: end")
}
