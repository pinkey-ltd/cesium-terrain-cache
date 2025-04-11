// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"context"
	"errors"
	"github.com/pinkey-ltd/cesium-terrain-cache/internal/adapter/handlers"
	"github.com/pinkey-ltd/cesium-terrain-cache/internal/adapter/store"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	// Get the tileset store
	store := &store.Store{}

	mux := http.NewServeMux()
	mux.HandleFunc("/{tileset}/layer.json", handlers.LayerHandler(ctx))
	mux.HandleFunc("/{tileset}/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.terrain?v={version}", handlers.TerrainHandler(store))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("Server is running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:" + err.Error())
	}
	// TODO: hotReload
	slog.Info("Server exiting")
}
