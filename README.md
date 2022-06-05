# featurez-api

> Backend api for featurez

## What Is Featurez?

Featurez is a feature flag management tool which stores your feature flags in Redis. It allows you to seamlessly create, delete, search, activate and deactivate your feature flag.

Link to the frontend:

## Getting Started

### Prerequisites

Go v1.18

### Required Services (Redis & PostgreSQL)

I recommend you run these services on docker! By running the command below you will deploy all the required services.

Run docker-compose file to deploy PostgreSQL to docker

```docker
cd .\featurez-api\

docker compose up -d
```

You can bring down all services with the following command. Keep in mind, Redis is meant for caching, you will lose all data once Redis has restarted.

```docker
docker compose down
```

### Accessing DB

Once you start the application, Gorm will auto migrate the schema to the database. You can access the database with the below commands.

```sh
docker exec -it featurez-api-db-1 sh

psql -U gorm

\d settings
```

### Run the App

Most environment variables are given defaults, but if you would like a different config, reference .env.examples and copy the environment variables into a .env file.

Run command

```sh
go run main.go
```

TODO

1. Add password capabilities for Redis.

```sh
command: redis-server --save 20 1 --loglevel warning --requirepass <password>
```
