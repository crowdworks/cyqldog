package cyqldog

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Monitor represents a main process of monitoring.
type Monitor struct {
	configPath string
}

// NewMonitor returns an instance of Monitor.
func NewMonitor(configPath string) *Monitor {
	return &Monitor{
		configPath: configPath,
	}
}

// Run is a main routine of cyqldog.
func (m *Monitor) Run() error {
	// Load the configuration file.
	log.Printf("monitor: load config file: %s", m.configPath)
	config, err := newConfig(m.configPath)
	if err != nil {
		return err
	}

	// Connect to the database.
	ds, err := newDB(config.DB)
	if err != nil {
		return err
	}
	defer ds.Close()

	// Initialize notifiers.
	notifiers, err := newNotifiers(config.Notifiers)
	if err != nil {
		return err
	}

	// Make a task queue for monitoring job.
	q := make(chan Rule)

	// Make a scheduler for each rule.
	for i, rule := range config.Rules {
		scheduler := newScheduler(i, rule)
		go scheduler.run(q)
	}

	// Make a monitoring worker.
	// In order to limit the number of DB connection to 1 for monitoring,
	// only one worker should run.
	c := newChecker(ds, notifiers)
	go c.run(q)

	// Trap signals from OS for normal termination.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// The main goroutine keep on waiting here.
loop:
	for {
		select {
		case <-stop:
			log.Println("monitor: stopping")
			break loop
		}
	}

	return nil
}
