// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"flag"
	"fmt"
	"github.com/pinkey-ltd/cesium-terrain-server/handlers"
	"github.com/pinkey-ltd/cesium-terrain-server/internal/cache"
	"github.com/pinkey-ltd/cesium-terrain-server/stores/fs"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := flag.Uint("port", 8000, "the port on which the server listens")
	tilesetRoot := flag.String("dir", ".", "the root directory under which tileset directories reside")
	cacheable := flag.Bool("cacheable", false, "(optional) enable caching tiles")
	baseTerrainUrl := flag.String("base-terrain-url", "/tilesets", "base url prefix under which all tilesets are served")
	limit := NewLimitOpt()
	limit.Set("1MB")
	flag.Var(limit, "cache-limit", `the memory size in bytes beyond which resources are not cached. Other memory units can be specified by suffixing the number with KB, MB, GB or TB`)
	flag.Parse()

	// Get the tileset store
	store := fs.New(*tilesetRoot)

	r := http.NewServeMux()
	r.HandleFunc(*baseTerrainUrl+"/{tileset}/layer.json", handlers.LayerHandler(store))
	r.HandleFunc(*baseTerrainUrl+"/{tileset}/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.terrain", handlers.TerrainHandler(store))

	handler := handlers.AddCorsHeader(r)

	if *cacheable {
		slog.Debug(fmt.Sprintf("cacheable enabled for all resources: %s", *limit))
		handler = cache.NewCache(handler, limit.Value, handlers.NewLimit)
	}

	http.Handle("/", handler)

	slog.Info(fmt.Sprintf("server listening on port %d", *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		slog.Error(fmt.Sprintf("server failed: %s", err))
		os.Exit(1)
	}
}
