package main

import (
	"os"

	"github.com/vadim8q258475/geo-microservice/app"
	"github.com/vadim8q258475/geo-microservice/service"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	godotenv.Load(".env")

	port := os.Getenv("PORT")
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	service := service.NewGeoGrpcService(apiKey, secretKey, logger)

	server := grpc.NewServer()

	app := app.NewApp(service, server, logger, port)

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
