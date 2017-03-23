/*
Log of output:

2017/03/22 02:33:45 main.go:47: Starting server, date and time is 2017-03-22T02:33:45Z
2017/03/22 02:33:52 main.go:77: Stop signal received at 2017-03-22T02:33:52Z
2017/03/22 02:33:52 main.go:80: Shutting down server...
2017/03/22 02:33:52 main.go:85: Server shutdown successfully
2017/03/22 02:33:52 main.go:88: Shutdown complete
2017/03/22 02:33:52 main.go:26: HTTP error: http: Server closed

*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bouk/httprouter"
)

func main() {
	runtime.GOMAXPROCS(2)

	fmt.Println("Starting server...")

	router := httprouter.New()
	// Set routes here.

	// Create stop channel to watch for interrupts.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Create a new HTTP server.
	httpServer := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	// Start the server.
	ch := make(chan error)
	go func() {
		ch <- httpServer.ListenAndServe()
	}()

	// Wait for a signal or ListenAndServe to return an (unexpected) error.
	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Printf("Server shutdown error: %s\n", err)
		} else {
			logger.Printf("Server shutdown successfully")
		}
		// Wait for ListenAndServe to return, ignore the error.
		<-ch
	case err := <-ch:
		logger.Fatal(err)
	}
	logger.Println("Shutdown complete")
}
