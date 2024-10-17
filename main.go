package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trackerApp/configs"
	"trackerApp/docs"
	"trackerApp/internal/handlers"
	"trackerApp/internal/services"
	"trackerApp/pkg/httpServer"
	"trackerApp/pkg/postgres"

	"github.com/spf13/viper"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server TaskTracker server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "taskTracker.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if err := configs.Init(); err != nil {
		panic(err)
	}
	db, err := postgres.NewPostgresDb(postgres.PostgresConfig{
		Username:     viper.GetString("dbDevelopment.username"),
		Password:     viper.GetString("dbDevelopment.password"),
		Host:         viper.GetString("dbDevelopment.host"),
		Port:         viper.GetString("dbDevelopment.port"),
		DatabaseName: viper.GetString("dbDevelopment.dbname"),
		SslMode:      viper.GetString("dbDevelopment.sslmode"),
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	service := services.NewService(db)
	handler := handlers.NewHandler(service)

	server := new(httpServer.Server)
	go func() {
		if err := server.ListenAndServe(":8080", handler.InitRoutes()); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("shutting down...")
	}
	log.Println("Server exiting")
}
