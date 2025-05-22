package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yash3004/config_server/cmd"
	"github.com/yash3004/config_server/configurations"
	"github.com/yash3004/config_server/internal/transport/grpc_transport"
	"github.com/yash3004/config_server/internal/transport/http_transport"
	"github.com/yash3004/config_server/users"
)

func main() {
	cfg := cmd.GetConfigurations()
	mongoURI := flag.String("mongoURI", cfg.MongoURI, "MongoDB URI")
	grpcAddr := flag.String("grpcAddr", fmt.Sprintf(":%d", cfg.Bind.GRPC), "gRPC server address")
	httpAddr := flag.String("httpAddr", fmt.Sprintf(":%d", cfg.Bind.HTTP), "HTTP server address")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Received termination signal, shutting down...")
		cancel()
	}()

	db, err := configurations.InitMongoDB(ctx, *mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	userManager := users.NewUserManager(db)
	configManager := configurations.NewConfigManager(db, cfg.UseFile, "configs")

	grpcServer := grpc_transport.NewServer(userManager, configManager)
	go func() {
		log.Printf("Starting gRPC server on %s", *grpcAddr)
		if err := grpc_transport.StartGRPCServer(grpcServer, *grpcAddr); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	httpServer := http_transport.NewServer(userManager, configManager)
	go func() {
		log.Printf("Starting HTTP server on %s", *httpAddr)
		if err := httpServer.StartHTTPServer(*httpAddr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Server shutdown complete")
}
