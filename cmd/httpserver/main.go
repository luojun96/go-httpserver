package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/luojun96/go-httpserver/metric"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server ... ")
	metric.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/metrics", promhttp.Handler())
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter root handler.")
	timer := metric.NewTimer()
	defer timer.Observe()
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}

	io.WriteString(w, "Headers:\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}
