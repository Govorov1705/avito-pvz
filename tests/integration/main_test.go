package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Govorov1705/avito-pvz/configs"
	"github.com/Govorov1705/avito-pvz/internal/handlers"
	"github.com/Govorov1705/avito-pvz/internal/services"
	"github.com/Govorov1705/avito-pvz/internal/storages/postgresql"
	"github.com/Govorov1705/avito-pvz/internal/storages/postgresql/repositories"
	"github.com/Govorov1705/avito-pvz/pkg/logger"
	"github.com/Govorov1705/avito-pvz/pkg/middleware"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresContainer *postgres.PostgresContainer
	storage           *postgresql.Storage
	router            *gin.Engine
	userService       *services.UserService
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	postgresContainer, err = postgres.Run(ctx,
		"postgres:17.2",
		postgres.WithDatabase("tests"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(fmt.Sprintf("error creating container: %s\n", err.Error()))
	}

	connStr, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		panic(fmt.Sprintf("error getting connection string: %s\n", err.Error()))
	}

	configs.Cfg.Mode = "dev"
	configs.Cfg.DBURL = connStr

	logger.InitLogger()

	storage = postgresql.NewStorage()

	userRepo := repositories.NewUserRepository(storage.Pool)
	pvzRepo := repositories.NewPvzRepository(storage.Pool)
	productRepo := repositories.NewProductRepository(storage.Pool)
	receptionRepo := repositories.NewReceptionRepository(storage.Pool)

	productService := services.NewProductService(productRepo, receptionRepo, storage.Pool)
	receptionService := services.NewReceptionService(receptionRepo, storage.Pool)
	pvzService := services.NewPvzService(pvzRepo, receptionRepo, productRepo)
	userService = services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(userService)
	pvzHandler := handlers.NewPvzHandler(pvzService, receptionService, productService)
	receptionHandler := handlers.NewReceptionHandler(receptionService)
	productHandler := handlers.NewProductHandler(productService)

	router = gin.Default()
	router.POST("/dummyLogin", authHandler.DummyLogin)

	authenticated := router.Group("/")
	authenticated.Use(middleware.Auth)

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
	employeeOnly.POST("/products", productHandler.AddProductToOpenReception)

	code := m.Run()

	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		panic(fmt.Sprintf("error terminating container: %s\n", err.Error()))
	}

	os.Exit(code)
}
