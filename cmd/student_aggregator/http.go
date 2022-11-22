/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package studentaggregator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/internal/transport/handlers"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/factory"
	"github.com/abdukhashimov/student_aggregator/pkg/minio"
	"github.com/abdukhashimov/student_aggregator/pkg/mongodb"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		portNum, _ := cmd.Flags().GetInt("port")

		serveHttp(portNum)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	httpCmd.PersistentFlags().Int("port", 8080, "http server port to be served")
}

func serveHttp(port int) {
	cfg := config.Load(config.TRANSPORT_HTTP)
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

	storageClient := minio.NewClient(cfg.Storage)
	logger.Log.Info("Minio connection success")

	server := handlers.NewServer(db, storageClient, cfg)

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
