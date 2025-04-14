package logger

import (
	"fmt"
	"log"

	"github.com/Govorov1705/avito-pvz/configs"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	var err error

	switch configs.Cfg.Mode {
	case configs.ModeDev:
		Logger, err = zap.NewDevelopment()
	case configs.ModeProd:
		Logger, err = zap.NewProduction()
	default:
		Logger = zap.NewNop()
		fmt.Println("MODE env not set, using default no-op logger for testing")
	}

	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	Logger.Info("Logger initialized")
}
