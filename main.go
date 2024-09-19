package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var (
		verbose     bool
		logFile     string
		bindAddress string
		certFile    string
		keyFile     string
		endpoints   Endpoints
		mu          sync.Mutex
		app         App
	)

	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&logFile, "f", "", "Log file (default stdout)")
	flag.StringVar(&bindAddress, "b", "0.0.0.0:80", "Local address to listen on")
	flag.StringVar(&certFile, "crt", "", "Crt file for TLS")
	flag.StringVar(&keyFile, "key", "", "Key file for TLS")
	flag.Var(&endpoints, "a", "Endpoint in format 'url:host:port' (e.g. /ssh:localhost:22)")
	flag.BoolVar(&app.quiet, "q", false, "Do not print anything")
	flag.DurationVar(&app.timeout, "t", 2*time.Second, "Timeout for new connections")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	logger, err := NewLogger(logFile, level)
	if err != nil {
		app.PrintError(err.Error())
		return
	}
	slog.SetDefault(logger)

	for _, ep := range endpoints {
		http.HandleFunc(fmt.Sprintf("GET %s", ep.url), func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			requestIP := strings.SplitN(r.RemoteAddr, ":", 2)[0]
			var ln net.Listener
			ln, err = app.StartListener(requestIP, ep.endpoint)
			if err != nil {
				app.PrintError(err.Error())
			}
			addr, _ := ln.Addr().(*net.TCPAddr)
			port := fmt.Sprintf("%d", addr.Port)
			if _, err = io.WriteString(w, port); err != nil {
				app.LogError(err.Error())
			}
		})
	}

	logRequest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr, _ := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		app.LogInfo("http",
			"remoteAddr", addr,
			"agent", r.UserAgent(),
			"method", r.Method,
			"url", r.URL.String())
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	addr, _ := net.ResolveTCPAddr("tcp", bindAddress)
	if certFile != "" && keyFile != "" {
		if addr.Port == 80 {
			addr.Port = 443
		}
		app.LogInfo("http/listening", "addr", addr)
		err = http.ListenAndServeTLS(addr.String(), certFile, keyFile, logRequest)
	} else {
		app.LogInfo("Cert and key files are not provided: TLS is disabled...")
		app.LogInfo("http/listening", "addr", addr)
		err = http.ListenAndServe(addr.String(), logRequest)
	}
	if err != nil {
		app.PrintError(err.Error())
	}
}

func NewLogger(file string, level slog.Level) (*slog.Logger, error) {
	w := os.Stdout
	if file != "" {
		var err error
		w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})), nil
}
