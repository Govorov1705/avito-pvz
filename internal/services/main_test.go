package services_test

import (
	"os"
	"testing"

	"github.com/Govorov1705/avito-pvz/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	code := m.Run()
	os.Exit(code)
}
