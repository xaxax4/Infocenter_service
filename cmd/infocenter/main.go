package main

import (
	"Infocenter/internal/api/handlers"
	"Infocenter/internal/api/middlewares"
	"Infocenter/internal/api/router"
	"Infocenter/internal/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000" // default
	} else if port[0] != ':' {
		port = ":" + port
	}

	shutdownStr := os.Getenv("SERVER_SHUTDOWN_TIMEOUT")
	if shutdownStr == "" {
		shutdownStr = "5" // default
	}
	shutdownSeconds, err := strconv.Atoi(shutdownStr)
	if err != nil {
		shutdownSeconds = 5 // default
	}

	serverShutdownTimeout := time.Duration(shutdownSeconds) * time.Second

	server := http.Server{
		Addr:    port,
		Handler: middlewares.CORSMiddleware(router.MainRouter()),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// create middleman and start it to handle topics and clients
	middleman := services.NewMiddleMan()
	handlers.StartMiddleman(middleman)

	go func() {
		log.Printf("Server started on port: %s", port)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Println("Server shutdown timed out but some connections were still active!")
		} else {
			log.Fatalf("Error shutting down server: %v", err)
		}
	}

	log.Println("Server gracefully stopped")
}
