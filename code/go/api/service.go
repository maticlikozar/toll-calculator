package api

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"toll/api/audit"
	"toll/api/auth"
	"toll/api/config"
	httpHandler "toll/api/handler/http"
	"toll/api/restapi"
	"toll/api/service"

	"toll/internal/database"
	"toll/internal/log"
	"toll/internal/request"
	"toll/internal/sigctx"
)

var (
	flags *config.AppFlags
)

func Init() error {
	flags = config.Get()

	if err := flags.Validate(); err != nil {
		return err
	}

	log.Info("initializing services")
	service.Init()

	log.Info("creating database client and connection")
	database.Get()

	rand.NewSource(time.Now().UnixNano())

	return nil
}

func Start() error {
	var deadline = sigctx.New()

	log.WithFields(log.Fields{"address": flags.Svc.Addr}).Info("starting API service")

	go StartListener()

	<-deadline.Done()
	log.Print("stopping API service")

	return nil
}

func StartListener() {
	authorization := auth.Get()

	handler, err := restapi.NewServer(
		httpHandler.NewApiHandlers(),
		authorization,
		restapi.WithPathPrefix("/api/v1"),
		restapi.WithMiddleware(audit.Middleware),
	)
	if err != nil {
		log.Fatale(err, "error creating handler")
	}

	h := request.Logger(handler)
	h = request.RequestId(h)

	mux := http.NewServeMux()
	mux.Handle("/", h)

	srv := &http.Server{
		Addr:         flags.Svc.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	log.Infof("Listening on %s%s", flags.Svc.Domain, flags.Svc.Addr)

	// Listen to OS signals for program termination.
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		//nolint
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatale(err)
		}
	}()

	log.Info("Server started")

	<-osSignal

	log.Info("Server shutting down")

	// Timeout to prevent server.Shutdown taking forever waiting for connections to close.
	contextTimeout := 5 * time.Second
	ctxShutdown, cancelCtxShutdown := context.WithTimeout(context.Background(), contextTimeout)

	defer cancelCtxShutdown()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Fatale(err, "error gracefully shutting down the server")
	}

	log.Info("Server exited properly")
}
