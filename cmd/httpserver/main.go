package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server ... ")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	fmt.Println("vim-go")
	server := http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen:L %s\n", err.Error())
		}
	}()
	glog.V(2).Info("Server started.")
	<-done

	glog.V(2).Info("Server Stopped.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed: %+v", err)
	}
	glog.V(2).Info("Server Exited Properly.")
}
