package main

import (
	"flag"
	"os"
	"time"
)

func main() {
	var logFile string
	var app = App{endpoints: Endpoints{}}

	flag.StringVar(&logFile, "f", "", "Log file (default stdout)")
	flag.StringVar(&app.bindAddress, "b", "0.0.0.0:80", "Local address to listen on")
	flag.StringVar(&app.certFile, "crt", "", "Crt file for TLS")
	flag.StringVar(&app.keyFile, "key", "", "Key file for TLS")
	flag.Var(&app.endpoints, "a", "Endpoint in format 'url:host:port' (e.g. /ssh:localhost:22)")
	flag.BoolVar(&app.quiet, "q", false, "Do not print anything")
	flag.BoolVar(&app.verbose, "v", false, "Verbose mode")
	flag.DurationVar(&app.timeout, "t", 2*time.Second, "Timeout for accepting connections")
	flag.BoolVar(&app.noTarpit, "ntp", false, "Do not tarpit wrong requests")
	flag.Parse()

	if err := app.InitSLogger(logFile, 0); err != nil {
		app.PrintError(err.Error())
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		app.PrintError(err.Error())
		os.Exit(1)
	}
}
