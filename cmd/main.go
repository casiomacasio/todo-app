package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application
// @host localhost:8000
// @BasePath /
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}
	db, err := database.NewPostgresDB(database.Config{
		Host:     viper.GetString("postgres.db.host"),
		Port:     viper.GetString("postgres.db.port"),
		Username: viper.GetString("postgres.db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("postgres.db.dbname"),
		SSLMode:  viper.GetString("postgres.db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error initializing db: %s", err.Error())
	}
	rdb, err := database.NewRedisClient(database.RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       viper.GetInt("redis.db"),
	})
	if err != nil {
		logrus.Fatalf("failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	repos := repository.NewRepository(db)
	services := service.NewService(repos, rdb)
	handlers := handler.NewHandler(services, rdb)

	srv := new(server.Server)
	go func () {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %v", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured when shutting down the server %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured when closing db connection %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
