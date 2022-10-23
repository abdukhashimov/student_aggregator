package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/clitools"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/internal/transport/handlers"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/factory"
	"github.com/abdukhashimov/student_aggregator/pkg/mongodb"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no subcommand")
		return
	}

	comm := clitools.PasreCommand()

	switch comm.Action {
	case "http":
		fmt.Printf("HTTP Port is %d \n", comm.Port)
		runHTTP(comm.Port)
	case "grpc":
		fmt.Printf("GRPC Port is %d \n", comm.Port)
	default:
		fmt.Println("Unknown subcommand: " + os.Args[1])
	}
}

func runHTTP(port int) {
	cfg := config.Load()
	log, err := factory.Build(&cfg.Logging)
	if err != nil {
		panic(err)
	}
	cfg.Http.Port = port
	logger.SetLogger(log)
	log.Info("logger successfully initialized")

	mongoClient, err := mongodb.NewClient(cfg.MongoDB.URI, cfg.MongoDB.User, cfg.MongoDB.Password)
	if err != nil {
		panic(err)
	}

	db := mongoClient.Database(cfg.MongoDB.Database)
	log.Info("mongo db client successfully initialized")

	server := handlers.NewServer(db, cfg)

	go func() {
		err = server.Run(fmt.Sprintf("%d", cfg.Http.Port))
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Project.GracefulTimeoutSeconds))
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Info("shutting down")

		server.Shutdown(ctx)

		log.Info("shutdown successfully called")

		wg.Done()
	}(&wg)

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
	os.Exit(0)
}
