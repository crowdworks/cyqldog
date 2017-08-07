package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/crowdworks/cyqldog/monitor"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("main: start")

	var configPath string
	flag.StringVar(&configPath, "C", "./config/local.yml", "path to config file")
	flag.Parse()

	run(configPath)
	log.Println("main: end")
}

func run(configPath string) {
	log.Printf("main: load config file: %s", configPath)
	config, err := monitor.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := monitor.NewDB(config.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	statsd, err := monitor.NewStatsd(config.Statd)
	if err != nil {
		log.Fatal(err)
	}

	q := make(chan monitor.Rule)

	for i, rule := range config.Rules {
		scheduler := monitor.NewScheduler(i, rule)
		go scheduler.Run(q)
	}

	m := monitor.NewMonitor(db, statsd)
	go m.Run(q)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-stop:
			log.Println("main: stopping")
			break loop
		}
	}
}
