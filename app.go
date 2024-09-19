package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

type Item struct {
	listener net.Listener
	count    int
}

type App struct {
	quiet     bool
	timeout   time.Duration
	listeners map[string]*Item
}

func (app *App) PrintError(format string, args ...any) {
	if !app.quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func (app *App) LogInfo(format string, args ...any) {
	if !app.quiet {
		slog.Info(format, args...)
	}
}

func (app *App) LogDebug(format string, args ...any) {
	if !app.quiet {
		slog.Debug(format, args...)
	}
}

func (app *App) LogError(format string, args ...any) {
	if !app.quiet {
		slog.Error(format, args...)
	}
}

func (app *App) StartListener(clientIP, endpoint string) (net.Listener, error) {
	if app.listeners == nil {
		app.listeners = make(map[string]*Item)
	}
	item, ok := app.listeners[clientIP]
	if !ok {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, err
		}
		item = &Item{listener: ln, count: 0}
		app.listeners[clientIP] = item
		app.LogInfo("tcp/listening", "addr", ln.Addr())
		go app.AcceptLoop(clientIP, ln, endpoint)
	}
	return item.listener, nil
}

func (app *App) AcceptLoop(requestIP string, listener net.Listener, endpoint string) {
	ln, _ := listener.(*net.TCPListener)
	defer func() {
		_ = ln.Close()
		delete(app.listeners, requestIP)
		app.LogInfo("closed", "addr", ln.Addr())
	}()

	for {
		if err := ln.SetDeadline(time.Now().Add(app.timeout)); err != nil {
			app.LogDebug(err.Error())
			return
		}
		remoteConn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if item, ok := app.listeners[requestIP]; ok && item.count > 0 {
					continue
				}
			}
			app.LogError(err.Error())
			return
		}
		app.LogInfo("connected", "laddr", remoteConn.LocalAddr(), "raddr", remoteConn.RemoteAddr())
		addr, _ := remoteConn.RemoteAddr().(*net.TCPAddr)
		remoteIP := fmt.Sprintf("%s", addr.IP)
		item, ok := app.listeners[remoteIP]
		if !ok {
			app.LogError("ip_mismatch", "requestIP", requestIP, "clientIP", remoteIP)
			return
		}
		item.count += 1
		go func() {
			if err = app.Connect(remoteConn, endpoint); err != nil {
				app.LogError(err.Error())
			}
			item.count -= 1
		}()
	}
}

func (app *App) Connect(remoteConn net.Conn, localAddr string) error {
	defer func() {
		_ = remoteConn.Close()
		app.LogDebug("closed", "laddr", remoteConn.LocalAddr(), "raddr", remoteConn.RemoteAddr())
	}()

	localConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		return err
	}
	defer func() {
		_ = localConn.Close()
		app.LogDebug("closed", "laddr", localConn.LocalAddr(), "raddr", localConn.RemoteAddr())
	}()

	app.LogInfo("connected", "laddr", localConn.LocalAddr(), "raddr", localConn.RemoteAddr())
	go io.Copy(remoteConn, localConn)
	_, err = io.Copy(localConn, remoteConn)

	return err
}
