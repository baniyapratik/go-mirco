package sever

import (
	"context"
	"fmt"
	"log"
	mongo_driver "logger-service/mongo-driver"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v6"
)

type Server struct {
	cfg         Config
	errc        chan error
	mongoClient mongo_driver.Interactor
}

// New creates a new HTTP server and setup routing
func New() (*Server, error) {
	type cfgFunc func() error

	server := &Server{}
	var err error
	for _, step := range []cfgFunc{
		server.loadConfig,
		server.connectToMongo,
	} {
		if err = step(); err != nil {
			return nil, fmt.Errorf("failed to initialize step %w", err)
		}
	}
	return server, nil
}

func (s *Server) ListenAndServe() error {
	h := s.newHTTPHandler()
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.WebPort)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	go func() {
		s.errc <- httpServer.ListenAndServe()
	}()

	// setup interrupt handlers to ensure we stop gracefully
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	log.Println("listening:", addr)
	select {
	case sig := <-signals:
		log.Println("caught signal, shutting down. Signal:", sig.String())
	case err := <-s.errc:
		log.Println("server error, shutting down", err)
	}

	serverShutdownTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	// close connection
	defer func() {
		if err := s.mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("shutdown http server:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
	return nil
}

func (s *Server) loadConfig() error {
	if err := env.Parse(&s.cfg); err != nil {
		return fmt.Errorf("configuration could not be parsed for environment: %w", err)
	}
	return nil
}

func (s *Server) connectToMongo() error {
	var err error
	mongoDBUsername := ""
	mongoDBPassword := ""
	mongoUrl := s.cfg.MongoURL
	s.mongoClient, err = mongo_driver.New(
		mongo_driver.WithLogger(log.Logger{}),
		mongo_driver.WithMongoDBURL(mongoUrl),
		mongo_driver.WithUsername(mongoDBUsername),
		mongo_driver.WithPassword(mongoDBPassword),
		mongo_driver.WithDatabaseName(s.cfg.ServiceDatabaseName),
		mongo_driver.WithCollectionName(s.cfg.ServiceLogsCollectionName))
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb %w", err)
	}
	return nil
}
