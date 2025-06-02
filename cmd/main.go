package main

import (
	"os"

	"github.com/casiomacasio/todo-app/internal/handler"
	"github.com/casiomacasio/todo-app/internal/repository"
	"github.com/casiomacasio/todo-app/internal/server"
	"github.com/casiomacasio/todo-app/internal/service"
	"github.com/casiomacasio/todo-app/pkg/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}
	db, err := database.NewPostgresDB(database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error initializing db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %v", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
