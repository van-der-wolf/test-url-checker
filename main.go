/*

 */
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BrightLocal/test-url-checker-ms/checker"
	"github.com/BrightLocal/test-url-checker-ms/protocol"
)

const listenOn = "10007"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("/"+protocol.MethodCheckURLs, http.HandlerFunc(CheckURLs))

	server := &http.Server{
		Addr:    ":" + listenOn,
		Handler: mux,
	}

	go func() {
		<-signals
		log.Printf("Graceful shutdown initiated...")
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Printf("Listen on :%s", listenOn)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen error: %s\n", err)
	}
}

func CheckURLs(rw http.ResponseWriter, r *http.Request) {
	var req protocol.ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(rw).Encode(protocol.ProfileResponse{Error: err.Error()}); err != nil {
			log.Printf("failed to marshal error response: %s", err)
		}
		return
	}

	resp := &protocol.ProfileResponse{URLCodes: checker.Query(req.URLs)}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		log.Printf("failed to marshal response: %s", err)
	}
}
