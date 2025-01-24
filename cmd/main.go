package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/PriyanshuDubey123/Assignment_Kirana_Club/apiroutes"
	"github.com/PriyanshuDubey123/Assignment_Kirana_Club/internals/db"
	"github.com/PriyanshuDubey123/Assignment_Kirana_Club/internals/repository"
	"github.com/PriyanshuDubey123/Assignment_Kirana_Club/internals/store"
	"github.com/PriyanshuDubey123/Assignment_Kirana_Club/message_broker"
)

func init() {
	wd, _ := os.Getwd() // Get the current working directory
	log.Printf("Current working directory: %s", wd)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("%v", err)
	}
	db.DBConnection()
}

func main() {

	_, err := message_broker.NewProducer(os.Getenv("RABBITMQ_URL"), os.Getenv("QUEUE_NAME"))
	if err != nil {
		log.Printf("Error from Producer %v\n", err.Error())
	}
	consumer, err := message_broker.NewConsumer(os.Getenv("RABBITMQ_URL"), os.Getenv("QUEUE_NAME"))
	if err != nil {
		log.Printf("%v", err)
	}
	go func() {
		if err := consumer.Start(); err != nil {
			log.Printf("Error from  consumer %v\n", err)
		}
	}()
	store.NewStoreManager()
	repository.NewStoreService()
	router := gin.Default()
	apiroutes.StoreVisitServiceRoutes(router)
	srv := http.Server{
		Handler: router.Handler(),
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 2 milli seconds.")
	}
	log.Println("Server exiting")
}
