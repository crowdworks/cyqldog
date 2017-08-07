package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/crowdworks/cyqldog/cyqldog"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("main: start")

	// Parse the argument's flag.
	var configPath string
	flag.StringVar(&configPath, "C", "./config/local.yml", "path to config file")
	flag.Parse()

	run(configPath)
	log.Println("main: end")
}

func run(configPath string) {
	// Load the configuration file.
	log.Printf("main: load config file: %s", configPath)
	config, err := cyqldog.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database.
	db, err := cyqldog.NewDB(config.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connect to the dogstatsd.
	statsd, err := cyqldog.NewStatsd(config.Statd)
	if err != nil {
		log.Fatal(err)
	}

	// Make a task queue for monitoring job.
	q := make(chan cyqldog.Rule)

	// Make a scheduler for each rule.
	for i, rule := range config.Rules {
		scheduler := cyqldog.NewScheduler(i, rule)
		go scheduler.Run(q)
	}

	// Make a monitoring worker.
	// In order to limit the number of DB connection to 1 for monitoring,
	// only one worker should run.
	m := cyqldog.NewMonitor(db, statsd)
	go m.Run(q)

	// Trap signals from OS for normal termination.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// The main goroutine keep on waiting here.
loop:
	for {
		select {
		case <-stop:
			log.Println("main: stopping")
			break loop
		}
	}
}
