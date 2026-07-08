package main

import (
	"context"
	"errors"
	"fmt"
	"go-live-server/cli"
	"go-live-server/reload"
	"go-live-server/server"
	"go-live-server/watcher"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func formatAddress(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}

const (
	debounceInternval = 100 * time.Millisecond
	shutdownTimeout   = 5 * time.Second
)

func mustLoadConfig() cli.Config {
	config := cli.Parse()

	if err := cli.ValidateConfig(config); err != nil {
		log.Fatal(err)
	}

	configNormalized, err := cli.NormalizeConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	return configNormalized
}

func mustCreateWatcher(
	root string,
) *watcher.Watcher {
	w, err := watcher.NewWatcher(root)

	if err != nil {
		log.Fatal(err)
	}

	ignoreMatcher := watcher.NewIgnoreMatcher()

	err = watcher.WalkAndWatch(w, root, ignoreMatcher)

	if err != nil {
		log.Fatal(err)
	}

	w.Start()

	return w
}

func newHTTPServer(root string, port int, reloader *reload.Reloader) *http.Server {
	router := server.NewRouter(root, reloader)

	return &http.Server{
		Addr:    formatAddress(port),
		Handler: router,
	}
}

func startHTTPServer(
	server *http.Server,
) {
	err := server.ListenAndServe()

	if err == nil {
		return
	}

	if errors.Is(err, http.ErrServerClosed) {
		return
	}

	log.Fatal(err)

}

func handleEvents(
	ctx context.Context,
	events <-chan watcher.Event,
	reloader *reload.Reloader,
) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-events:
			if !ok {
				return
			}

			fmt.Printf("changed: %s (%s)\n", event.Path, event.Type)

			reloader.Broadcast()
		}
	}
}

func shutdown(
	parent context.Context,
	server *http.Server,
	watcher *watcher.Watcher,
) {
	ctx, cancel := context.WithTimeout(
		parent, shutdownTimeout,
	)

	defer cancel()

	log.Println("closing watcher")

	if err := watcher.Close(); err != nil {
		log.Printf("watcher close error: %v\n", err)
	}

	log.Println("Shutting down http server")

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("http shutdown error: %v\n", err)
	}

	log.Println("shutdown complete")

}

func main() {
	config := mustLoadConfig()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	w := mustCreateWatcher(config.Directory)

	events := watcher.Debounce(
		w.Events(),
		debounceInternval,
	)

	reloader := reload.New()

	go handleEvents(ctx, events, reloader)

	httpServer := newHTTPServer(
		config.Directory,
		config.Port,
		reloader,
	)

	go startHTTPServer(httpServer)

	fmt.Printf("server listening on http://localhost:%d\n", config.Port)

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdown(ctx, httpServer, w)
}
