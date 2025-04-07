package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type TilesetStatus byte

const (
	NOT_SUPPORTED TilesetStatus = iota
	NOT_FOUND
	FOUND
)

var ErrNoItem = errors.New("item not found")

type Store struct{}

// readFile reads the specified file and returns its contents as a byte slice or an error if the file is not accessible.
// If the file does not exist, it logs a debug message and returns ErrNoItem.
func (s *Store) readFile(filename string) ([]byte, error) {
	body, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug(fmt.Sprintf("not found file: " + filename))
			return nil, ErrNoItem
		}
		return nil, err
	}

	slog.Debug("file store: load: " + filename)
	return body, nil
}

// Tile Load a terrain tile on disk into the Terrain structure.
func (s *Store) Tile(tileset string, tile *Terrain) (err error) {
	filename := filepath.Join(
		tileset,
		strconv.FormatUint(tile.Z, 10),
		strconv.FormatUint(tile.X, 10),
		strconv.FormatUint(tile.Y, 10)+".terrain")

	body, err := s.readFile(filename)
	if err != nil {
		return
	}

	err = tile.UnmarshalBinary(body)
	return
}

func (s *Store) Layer(tileset string) ([]byte, error) {
	filename := filepath.Join(tileset, "layer.json")
	return s.readFile(filename)
}

func (s *Store) TilesetStatus(tileset string) (status TilesetStatus) {
	// check whether the tile directory exists
	_, err := os.Stat(filepath.Join(tileset))
	if err != nil {
		if os.IsNotExist(err) {
			return NOT_FOUND
		}
	}

	return FOUND
}
