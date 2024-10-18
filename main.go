package main

import (
	"flag"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"time"
)

type Config struct {
	TZ        string
	Bind      string
	CrtFile   string
	KeyFile   string
	Endpoints map[string]string
	Quiet     bool
	Verbose   bool
	Timeout   int
	NoTarpit  bool
}

var configPath = "/etc/concierge/concierge.toml"

var config = Config{
	TZ:      "Europe/Vilnius",
	Bind:    "0.0.0.0:80",
	CrtFile: "server.crt",
	KeyFile: "server.key",
	Timeout: 2,
}

func main() {
	flag.StringVar(&configPath, "c", configPath, "Path to TOML config file")
	flag.Parse()

	if data, err := os.ReadFile(configPath); err != nil {
		log.Fatal(err)
	} else if err = toml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	if loc, err := time.LoadLocation(config.TZ); err != nil {
		panic(err)
	} else {
		time.Local = loc
	}

	app := App{config: config}
	app.InitSLogger(0)

	if err := app.Run(); err != nil {
		app.PrintError(err.Error())
		os.Exit(1)
	}
}
