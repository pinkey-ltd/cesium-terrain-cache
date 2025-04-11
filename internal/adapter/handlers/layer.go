package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pinkey-ltd/cesium-terrain-cache/internal/adapter/store"
	"log/slog"
	"net/http"
)

// LayerHandler An HTTP handler which returns a tileset's `layer.json` file
func LayerHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err   error
			layer []byte
		)

		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				slog.Error(err.Error())
			}
		}()

		vars := mux.Vars(r)

		// Try and get a `layer.json` from the store
		layer, err = store.Layer(vars["tileset"])
		if err == store.ErrNoItem {
			err = nil // don't persist this error
			if store.TilesetStatus(vars["tileset"]) == store.NotFound {
				http.Error(w,
					fmt.Errorf("the tileset `%s` does not exist", vars["tileset"]).Error(),
					http.StatusNotFound)
				return
			}

			// the directory exists: send the default `layer.json`
			layer = []byte(`{
  "tilejson": "2.1.0",
  "format": "heightmap-1.0",
  "version": "1.0.0",
  "scheme": "tms",
  "tiles": ["{z}/{x}/{y}.terrain"]
}`)
		} else if err != nil {
			return
		}

		headers := w.Header()
		headers.Set("Content-Type", "application/json")
		w.Write(layer)
	}
}

// TerrainHandler An HTTP handler which returns a terrain tile resource
func TerrainHandler(store *store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			t   store.Terrain
			err error
		)

		// get the tile coordinate from the URL
		vars := mux.Vars(r)
		err = t.ParseCoord(vars["x"], vars["y"], vars["z"], vars["version"])
		if err != nil {
			return
		}

		// Try and get a tile from the store
		err = store.Tile(vars["tileset"], &t)
		if err == store.ErrNoItem {
			if store.TilesetStatus(vars["tileset"]) == store.NOT_FOUND {
				err = nil
				http.Error(w,
					fmt.Errorf("The tileset `%s` does not exist", vars["tileset"]).Error(),
					http.StatusNotFound)
				return
			}

			if t.IsRoot() {
				// serve up a blank tile as it is a missing root tile
				data, err := assets.Asset("data/smallterrain-blank.terrain")
				if err != nil {
					return
				} else {
					err = t.UnmarshalBinary(data)
					if err != nil {
						return
					}
				}
			} else {
				err = nil
				http.Error(w, errors.New("The terrain tile does not exist").Error(), http.StatusNotFound)
				return
			}
		} else if err != nil {
			return
		}

		body, err := t.MarshalBinary()
		if err != nil {
			return
		}

		// send the tile to the client
		headers := w.Header()
		headers.Set("Content-Type", "application/octet-stream")
		headers.Set("Content-Encoding", "gzip")
		headers.Set("Content-Disposition", "attachment;filename="+vars["y"]+".terrain")
		w.Write(body)
	}
}
