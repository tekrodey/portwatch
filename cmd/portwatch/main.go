package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}

	alerter := alert.New(os.Stdout)
	mon := monitor.NewMonitor(cfg)

	fmt.Printf("portwatch: starting — scanning ports %d-%d every %s\n",
		cfg.StartPort, cfg.EndPort, cfg.Interval)

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			changes, err := mon.Scan()
			if err != nil {
				fmt.Fprintf(os.Stderr, "portwatch: scan error: %v\n", err)
				continue
			}
			if len(changes) > 0 {
				alerter.Notify(changes)
			}
		case sig := <-sigs:
			fmt.Printf("\nportwatch: received %s, shutting down\n", sig)
			return
		}
	}
}
