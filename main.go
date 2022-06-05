package main

import (
	"context"
	"featurez/api/feature"
	"featurez/api/settings"
	"featurez/clients"
	"featurez/config"
	"featurez/models"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// process environment variables
	err = envconfig.Process("", &config.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	// load postgresql
	postgresDNS := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Cfg.PostgresHost,
		config.Cfg.PostgresPort,
		config.Cfg.PostgresUser,
		config.Cfg.PostgresPassword,
		config.Cfg.PostgresDBName,
	)
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

	defaultRedisAddress := config.Cfg.RedisHost + ":" + config.Cfg.RedisPort
	if config.Cfg.RedisNoDB {
		clients.Redis = clients.NewRedisClient(defaultRedisAddress)
	} else {
		clients.Redis = clients.NewRedisClient(usrSettings.RedisAddress)
	}

	pong, err := clients.Redis.Client.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis Error:", err)
	}

	log.Printf("%s! Redis successfully connected", pong)

	mux := http.NewServeMux()

	// routes
	feature.Routes(config.Cfg.APIRoute, mux, "/feature")
	settings.Routes(config.Cfg.APIRoute, mux, "/settings")

	address := config.Cfg.FeaturezHost + ":" + config.Cfg.FeaturezPort

	log.Printf("Featurez-api is alive and listening on %s", address)

	http.ListenAndServe(address, mux)
}
