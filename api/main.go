package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	. "github.com/AbhishekKaushik16/instagram/api/config"
	"github.com/AbhishekKaushik16/instagram/api/routers"
	log "github.com/sirupsen/logrus"
)

var config Config

func init() {
	config.Read()
}

func main() {
	// Start the HTTP server

	log.SetOutput(os.Stdout)
	logFormatter := new(LogFormatter)
	logFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logFormatter.LevelDesc = []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
	log.SetFormatter(logFormatter)

	r := routers.Routers()

	srv := &http.Server{
		Handler:      r,
		Addr:         config.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	fmt.Println("Press Ctrl+C to shutdown the server")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Server shutting down...")
	err := srv.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("Error occurred during server shutdown: %v\n", err)
	}
	fmt.Println("Server gracefully stopped")
}
