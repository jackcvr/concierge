package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const ContentLength = 1024 * 64 * 10

type App struct {
	Logger
	bindAddress string
	certFile    string
	keyFile     string
	timeout     time.Duration
	endpoints   Endpoints
	noTarpit    bool
	mu          sync.Mutex
}

func (app *App) Run() error {
	for url, ep := range app.endpoints {
		http.HandleFunc(fmt.Sprintf("GET %s", url), func(w http.ResponseWriter, r *http.Request) {
			app.mu.Lock()
			defer app.mu.Unlock()

			requestIP := strings.SplitN(r.RemoteAddr, ":", 2)[0]
			ln, err := app.StartListener(requestIP, ep)
			if err != nil {
				app.PrintError(err.Error())
			}
			addr, _ := ln.Addr().(*net.TCPAddr)
			port := fmt.Sprintf("%d", addr.Port)
			if _, err = w.Write([]byte(port)); err != nil {
				app.Error(err.Error())
			}
		})
	}

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content := []byte("User-agent: *\nDisallow: /\n")
		if _, err := w.Write(content); err != nil {
			app.Error(err.Error())
		}
	})

	if !app.noTarpit {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", ContentLength))
			w.WriteHeader(http.StatusNotFound)
			payload := []byte("💔")
			for {
				time.Sleep(time.Second)
				if _, err := fmt.Fprint(w, payload); err != nil {
					app.Debug(err.Error())
					return
				}
				w.(http.Flusher).Flush()
			}
		})
	}

	logRequest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr, _ := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		start := time.Now()
		defer func() {
			lifetime := int(time.Now().Sub(start).Seconds())
			app.Info("http/closed", "remoteAddr", addr, "url", r.URL.String(), "lifetime", lifetime)
		}()
		app.Info("http/connected",
			"remoteAddr", addr,
			"agent", r.UserAgent(),
			"method", r.Method,
			"url", r.URL.String())
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	addr, _ := net.ResolveTCPAddr("tcp", app.bindAddress)
	if app.certFile != "" && app.keyFile != "" {
		if addr.Port == 80 {
			addr.Port = 443
		}
		app.Info("http/listening", "addr", addr)
		return http.ListenAndServeTLS(addr.String(), app.certFile, app.keyFile, logRequest)
	} else {
		app.Info("Cert and key files are not provided: TLS is disabled...")
		app.Info("http/listening", "addr", addr)
		return http.ListenAndServe(addr.String(), logRequest)
	}
}

func (app *App) StartListener(clientIP, endpoint string) (net.Listener, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	app.Info("tcp/listening", "addr", ln.Addr())
	go app.Accept(clientIP, ln, endpoint)
	return ln, nil
}

func (app *App) Accept(requestIP string, listener net.Listener, endpoint string) {
	ln, _ := listener.(*net.TCPListener)
	defer func() {
		_ = ln.Close()
		app.Info("tcp/closed", "addr", ln.Addr())
	}()

	if err := ln.SetDeadline(time.Now().Add(app.timeout)); err != nil {
		app.Error(err.Error())
		return
	}
	for {
		remoteConn, err := ln.Accept()
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				app.Error(err.Error())
			}
			return
		}
		app.Info("tcp/connected", "laddr", remoteConn.LocalAddr(), "raddr", remoteConn.RemoteAddr())
		addr, _ := remoteConn.RemoteAddr().(*net.TCPAddr)
		remoteIP := fmt.Sprintf("%s", addr.IP)
		if requestIP != remoteIP {
			_ = remoteConn.Close()
			app.Error("ip_mismatch", "requestIP", requestIP, "clientIP", remoteIP)
			continue
		}
		go app.Connect(remoteConn, endpoint)
		return
	}
}

func (app *App) Connect(remoteConn net.Conn, localAddr string) {
	defer func() {
		_ = remoteConn.Close()
		app.Debug("tcp/closed", "laddr", remoteConn.LocalAddr(), "raddr", remoteConn.RemoteAddr())
	}()

	localConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		app.Error(err.Error())
		return
	}
	defer func() {
		_ = localConn.Close()
		app.Debug("tcp/closed", "laddr", localConn.LocalAddr(), "raddr", localConn.RemoteAddr())
	}()

	app.Info("tcp/connected", "laddr", localConn.LocalAddr(), "raddr", localConn.RemoteAddr())
	go io.Copy(remoteConn, localConn)
	if _, err = io.Copy(localConn, remoteConn); err != nil {
		app.Error(err.Error())
	}
}
