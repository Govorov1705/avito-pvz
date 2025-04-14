package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Govorov1705/avito-pvz/configs"
	server "github.com/Govorov1705/avito-pvz/internal/grpc_server"
	"github.com/Govorov1705/avito-pvz/internal/handlers"
	"github.com/Govorov1705/avito-pvz/internal/metrics"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/internal/storages/postgresql"
	"github.com/Govorov1705/avito-pvz/internal/storages/postgresql/repositories"
	"github.com/Govorov1705/avito-pvz/pkg/logger"
	"github.com/Govorov1705/avito-pvz/pkg/middleware"
	pvz_v1 "github.com/Govorov1705/avito-pvz/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

func main() {
	// Different inits
	configs.InitConfig()
	logger.InitLogger()
	storage := postgresql.NewStorage()
	defer storage.Pool.Close()

	// Repos, services, handlers
	userRepository := repositories.NewUserRepository(storage.Pool)
	pvzRepository := repositories.NewPvzRepository(storage.Pool)
	receptionRepository := repositories.NewReceptionRepository(storage.Pool)
	productRepository := repositories.NewProductRepository(storage.Pool)

	userService := services.NewUserService(userRepository)
	pvzService := services.NewPvzService(pvzRepository, receptionRepository, productRepository)
	receptionService := services.NewReceptionService(receptionRepository, storage.Pool)
	productService := services.NewProductService(productRepository, receptionRepository, storage.Pool)

	authHandler := handlers.NewAuthHandler(userService)
	pvzHandler := handlers.NewPvzHandler(pvzService, receptionService, productService)
	receptionHandler := handlers.NewReceptionHandler(receptionService)
	productHandler := handlers.NewProductHandler(productService)

	// HTTP server
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	router.Use(middleware.Metrics)
	router.POST("/dummyLogin", authHandler.DummyLogin)
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	authenticated := router.Group("/")
	authenticated.Use(middleware.Auth)
	authenticated.GET("/pvz", pvzHandler.PvzInfo)

	moderatorOnly := authenticated.Group("/")
	moderatorOnly.Use(middleware.ModeratorOnly)
	moderatorOnly.POST("/pvz", pvzHandler.CreatePvz)

	employeeOnly := authenticated.Group("/")
	employeeOnly.Use(middleware.EmployeeOnly)
	employeeOnly.POST("/receptions", receptionHandler.CreateReception)
	employeeOnly.POST(
		"/pvz/:pvzId/close_last_reception",
		pvzHandler.CloseLastReception,
	)
	employeeOnly.POST(
		"/pvz/:pvzId/delete_last_product",
		pvzHandler.DeleteLastReceptionProduct,
	)
	employeeOnly.POST("/products", productHandler.AddProductToOpenReception)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("listen:", zap.Error(err))
		}
	}()

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		logger.Logger.Fatal("failed to listen:", zap.Error(err))
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pvz_v1.RegisterPVZServiceServer(grpcServer, server.NewPvzServer(pvzRepository))
	go func() {
		grpcServer.Serve(lis)
	}()

	// Metrics server
	metrics.StartMetricsServer()

	// Graceful shutdown
	// Listen for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutdown servers signal")

	// Create Contex with time for graceful shutdown for both servers
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpCancel()

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer grpcCancel()

	// Gracefully shutdown HTTP server
	logger.Logger.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(httpCtx); err != nil {
		logger.Logger.Error("Server Shutdown:", zap.Error(err))
	}

	// Gracefully shutdown gRPC server
	logger.Logger.Info("Shutting down gRPC server...")
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	// Force stop gRPC server if context timed out
	select {
	case <-grpcCtx.Done():
		grpcServer.Stop()
	case <-stopped:
	}

	logger.Logger.Info("All servers exited")
}
