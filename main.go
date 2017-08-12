package main

import (
	"flag"
	"log"

	"github.com/crowdworks/cyqldog/cyqldog"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("main: start")

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
