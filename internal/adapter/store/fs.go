package store

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

const pathPrefix = "data"

type TilesetStatus byte

const (
	StatusUnknown TilesetStatus = iota
	NotSupported
	NotFound
	Found
)

var ErrNoItem = errors.New("item not found")

type Store struct {
	tileset string
}

// readFile reads the specified file and returns its contents as a byte slice or an error if the file is not accessible.
// If the file does not exist, it logs a debug message and returns ErrNoItem.
func (s *Store) readFile(filename string) ([]byte, error) {
	pathRaw := filepath.Join(pathPrefix, filename)
	body, err := os.ReadFile(pathRaw)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug(fmt.Sprintf("not found file: " + filename))
			return nil, ErrNoItem
		}
		return nil, err
	}

	slog.Debug("file store load: " + filename)
	return body, nil
}

// Tile Load a terrain tile on disk into the Terrain structure.
func (s *Store) Tile(tile *Terrain) error {
	filename := filepath.Join(
		s.tileset,
		strconv.FormatUint(tile.Z, 10),
		strconv.FormatUint(tile.X, 10),
		strconv.FormatUint(tile.Y, 10)+".terrain")

	body, err := s.readFile(filename)
	if err != nil {
		return err
	}

	err = tile.UnmarshalBinary(body)
	return nil
}

// Layer retrieves the JSON-encoded metadata for the specified tileset by reading the "layer.json" file.
func (s *Store) Layer(tileset string) ([]byte, error) {
	filename := filepath.Join(tileset, "layer.json")
	return s.readFile(filename)
}

// Metadata retrieves the JSON-encoded metadata for the specified tileset by reading the "layer.json" file.
func (s *Store) Metadata(tileset string) ([]byte, error) {
	filename := filepath.Join(tileset, "metadata.json")
	return s.readFile(filename)
}

// TilesetStatus checks the existence and compatibility of the specified tileset directory and associated layer.json file.
// Returns NotFound if the tileset directory does not exist, NotSupported if "layer.json" file is missing, otherwise Found.
func (s *Store) TilesetStatus(tileset string) (status TilesetStatus) {
	// check whether the tile directory exists
	_, err := os.Stat(filepath.Join(tileset))
	if err != nil {
		if os.IsNotExist(err) {
			return NotFound
		} else {
			return StatusUnknown
		}
	}
	// check whether the layer.json exists
	_, err = os.Stat(filepath.Join(tileset, "layer.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return NotSupported
		} else {
			return StatusUnknown
		}
	}
	return Found
}
