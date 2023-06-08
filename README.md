
# Stark

Stark is an account service. Using golang as programming language, mysql as main database, and redis for cache.


## Run App
```sh
go run main.go
```

## Start Docker Compose Deployment
```sh
docker-compose -f docker-compose.yml up -d --build
```

## Stop Docker Compose Deployment
```sh
docker-compose -f docker-compose.yml down
```