package repository

import (
	"strconv"
)

// Terrain Representation of a terrain tile. This includes the x, y, z coordinate and
// the byte sequence of the tile itself. Note that terrain tiles are gzipped.
type Terrain struct {
	value   []byte
	X, Y, Z uint64
}

// MarshalBinary implements the encoding.MarshalBinary interface.
func (t *Terrain) MarshalBinary() ([]byte, error) {
	return t.value, nil
}

// UnmarshalBinary implements the encoding.UnmarshalBinary interface.
func (t *Terrain) UnmarshalBinary(data []byte) error {
	t.value = data
	return nil
}

// IsRoot returns true if the tile represents a root tile.
func (t *Terrain) IsRoot() bool {
	return t.Z == 0 &&
		(t.X == 0 || t.X == 1) &&
		t.Y == 0
}

// ParseCoord Parse x, y, z string coordinates and assign them to the Terrain
func (t *Terrain) ParseCoord(x, y, z, version string) error {
	xi, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return err
	}

	yi, err := strconv.ParseUint(y, 10, 64)
	if err != nil {
		return err
	}

	zi, err := strconv.ParseUint(z, 10, 64)
	if err != nil {
		return err
	}

	t.X = xi
	t.Y = yi
	t.Z = zi

	return nil
}
