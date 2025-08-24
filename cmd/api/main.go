package main

import (
	_ "WebProject/docs"
	"WebProject/internal/config"
	"WebProject/internal/repository/mongo"
	"WebProject/internal/service"
	"WebProject/internal/transport/rest/handler"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// @title Fiber Swagger API
// @version 2.0
// @description This is a sample server
// @termsOfService http://swagger.io/terms

// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	config.SetupSwagger(app)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoDatabase, err := config.SetupMongoDataBase(ctx, cancel)
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := mongo.NewUserRepository(mongoDatabase.Collection("users"))
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	userHandler.InitRoutes(app)

	port := viper.GetString("http.port")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}

func SetupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
