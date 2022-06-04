package main

import (
	"context"
	"featurez/api/feature"
	"featurez/api/settings"
	"featurez/clients"
	"featurez/models"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	FeaturezHost     string `split_words:"true" default:"localhost"`
	FeaturezPort     string `split_words:"true" default:"1000"`
	RedisNoDB        bool   `split_words:"true" default:"false"`
	RedisHost        string `split_words:"true" default:"localhost"`
	RedisPort        string `split_words:"true" default:"6379"`
	PostgresHost     string `split_words:"true" default:"localhost"`
	PostgresPort     string `split_words:"true" default:"5432"`
	PostgresUser     string `split_words:"true" default:"gorm"`
	PostgresPassword string `split_words:"true" default:"gorm"`
	PostgresDBName   string `split_words:"true" default:"gorm"`
	APIRoute         string `split_words:"true" default:"/api"`
}

var (
	cfg config
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// process environment variables
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// load postgresql
	postgresDNS := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName)
	clients.PostgresDB, err = clients.NewPostgresClient(postgresDNS)
	if err != nil {
		log.Fatal(err)
	}

	// migrate db schema
	err = clients.PostgresDB.Client.AutoMigrate(&models.Settings{})
	if err != nil {
		log.Fatal(err)
	}

	var usrSettings *models.Settings

	clients.PostgresDB.Client.First(&usrSettings)

	// insert a new settings row if none was found
	if usrSettings.ID == 0 {
		newSettings := &models.Settings{RedisAddress: ""}
		clients.PostgresDB.Client.Create(newSettings)
		clients.PostgresDB.Client.First(&usrSettings)
	}

	defaultRedisAddress := cfg.RedisHost + ":" + cfg.RedisPort
	if cfg.RedisNoDB {
		clients.Redis = clients.NewRedisClient(defaultRedisAddress)
	} else {
		clients.Redis = clients.NewRedisClient(usrSettings.RedisAddress)
	}

	pong, err := clients.Redis.Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis Error:", err)
	}

	log.Printf("%s! Redis successfully connected", pong)

	mux := http.NewServeMux()

	feature.Routes(cfg.APIRoute, mux, "/feature")
	settings.Routes(cfg.APIRoute, mux, "/settings")

	address := cfg.FeaturezHost + ":" + cfg.FeaturezPort

	log.Printf("Featurez-api is alive and listening on %s", address)

	http.ListenAndServe(address, mux)
}
