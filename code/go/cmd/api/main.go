package main

import (
	"time"

	"github.com/jnovack/flag"
	_ "github.com/joho/godotenv/autoload"

	service "toll/api"
	api "toll/api/config"

	"toll/internal/database"
	"toll/internal/errlog"
	"toll/internal/log"
)

func main() {
	flags("api",
		api.Flags,
		log.Flags,
		database.Flags,
	)

	log.SetTimeFormat(time.RFC3339)

	log.Info("initializing API service")

	if err := service.Init(); err != nil {
		log.WithFields(errlog.StackLog(err)).Fatale(err, "error initializing API service")
	}

	log.Info("starting API service")

	if err := service.Start(); err != nil {
		log.WithFields(errlog.StackLog(err)).Fatale(err, "error starting/running API service")
	}
}

func flags(prefix string, mountFlags ...func(...string)) {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")

	for _, mount := range mountFlags {
		mount(prefix)
	}

	flag.Parse()
}
