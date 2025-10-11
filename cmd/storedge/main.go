package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	httpadapter "stor-edge/internal/adapter/http"
	"sync"
	"time"
)

func main() {
	port := getenv("PORT", "8080")

	wt := sync.WaitGroup{}

	handler := httpadapter.NewRouter()

	srv := &http.Server{
		Addr: ":" + port,
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	wt.Add(1)
	go func() {
		defer wt.Done()
		log.Printf("starting StorEdge on :%s", port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
	}()
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}