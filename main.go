package main

import (
	"context"
	"featurez/api/feature"
	"featurez/api/settings"
	"featurez/config"
	"featurez/models"
	"featurez/services"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
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
	services.PostgresDB, err = services.NewPostgresClient(postgresDNS)
	if err != nil {
		log.Fatal(err)
	}

	// migrate db schema
	err = services.PostgresDB.Client.AutoMigrate(&models.Settings{})
	if err != nil {
		log.Fatal(err)
	}

	var usrSettings *models.Settings

	services.PostgresDB.Client.First(&usrSettings)

	// insert a new settings row if none was found
	if usrSettings.ID == 0 {
		newSettings := &models.Settings{RedisAddress: ""}
		services.PostgresDB.Client.Create(newSettings)
		services.PostgresDB.Client.First(&usrSettings)
	}

	defaultRedisAddress := config.Cfg.RedisHost + ":" + config.Cfg.RedisPort
	if config.Cfg.RedisNoDB {
		services.Redis = services.NewRedisService(defaultRedisAddress)
	} else {
		services.Redis = services.NewRedisService(usrSettings.RedisAddress)
	}

	pong, err := redis.Client.Ping(*services.Redis.Client, context.Background()).Result()
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
