package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Cfg struct {
	LoggerMode     string
	HTTPHost       string
	PostgreSQLDSN  string
	ForeignApiUrl  string
	TimeOut        int
	MigrationsPath string
}

const (
	loggerMode     = "LOGGER_MODE"
	httpHost       = "HTTP_HOST"
	postgresDsn    = "POSTGRESQL_DSN"
	foreignApiUrl  = "FOREIGN_API_URL"
	migrationsPath = "MIGRATIONS_PATH"
	timeout        = "TIMEOUT"
)

func NewConfig() (cfg Cfg) {
	err := godotenv.Load("././internal/config/.env")
	if err != nil {
		log.Println("cannot read config file\ndefault config accepted")
		return defaultConfig()
	}
	cfg.LoggerMode = os.Getenv(loggerMode)
	cfg.HTTPHost = os.Getenv(httpHost)
	cfg.PostgreSQLDSN = os.Getenv(postgresDsn)
	cfg.ForeignApiUrl = os.Getenv(foreignApiUrl)
	cfg.MigrationsPath = os.Getenv(migrationsPath)
	t, err := strconv.Atoi(os.Getenv(timeout))
	if err != nil {
		log.Println("cannot convert timeout")
	}
	cfg.TimeOut = t
	log.Println("config file reading success")
	return
}

func defaultConfig() (cfg Cfg) {
	cfg = Cfg{
		LoggerMode:     "debug",
		HTTPHost:       ":8080",
		PostgreSQLDSN:  "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		ForeignApiUrl:  "localhost",
		TimeOut:        10,
		MigrationsPath: "file://././internal/repository/postgresql/migrations",
	}
	return
}

var All = NewConfig()
