package core

// TileJSON See <https://github.com/mapbox/tilejson-spec>
type TileJSON string

const (
	VERSION1   TileJSON = "1.0.0"
	VERSION2   TileJSON = "2.0.0"
	VERSION201 TileJSON = "2.0.1"
	VERSION210 TileJSON = "2.1.0"
	VERSION220 TileJSON = "2.2.0"
	VERSION3   TileJSON = "3.0.0"
)

type Layer struct {
	terrains []*Terrain
	content  Content
}
type Terrain struct {
	value   []byte
	X, Y, Z uint64
}

type Content struct { // Version 2.1.0
	Tilejson    TileJSON // ALL VERSION
	Name        string   // ALL VERSION
	Description string
	Version     string
	Format      string
	Attribution string
	Schema      string
	Tiles       []string
	Projection  string
	Bounds      []float64
	Available   [][]struct {
		StartX int
		StartY int
		EndX   int
		EndY   int
	}
}

type LayerRepository interface {
	GetContent(tileset string) (*Content, error)
	GetTerrains(tileset string) ([]*Terrain, error)
	Save(tileset string, layer *Layer) error
	Delete(tileset string) error
	Create(path string) error
}

type ContentV3 struct {
	Name       string
	Type       string
	Properties struct {
		Tilejson struct {
			Type    string
			Pattern string
		}
		Tiles struct {
			Type  string
			Items struct {
				Type string
			}
		}
		VectorLayers struct {
			Type  string
			Items struct {
				Type       string
				Properties struct {
					Id struct {
						Type string
					}
					Fields struct {
						Type                 string
						AdditionalProperties struct {
							Type string
						}
					}
					Description struct {
						Type string
					}
					Maxzoom struct {
						Type string
					}
					Minzoom struct {
						Type string
					}
				}
				Required             []string
				AdditionalProperties bool
			}
		}
		Attribution struct {
			Type string
		}
		Bounds struct {
			Type  string
			Items struct {
				Type string
			}
		}
		Center struct {
			Type  string
			Items struct {
				Type string
			}
		}
		Data struct {
			Type  string
			Items struct {
				Type string
			}
		}
		Description struct {
			Type string
		}
		Fillzoom struct {
			Minimum int
			Maximum int
			Type    string
		}
		Grids struct {
			Type  string
			Items struct {
				Type string
			}
		}
		Legend struct {
			Type string
		}
		Maxzoom struct {
			Minimum int
			Maximum int
			Type    string
		}
		Minzoom struct {
			Minimum int
			Maximum int
			Type    string
		}
		Name struct {
			Type string
		}
		Scheme struct {
			Type string
		}
		Template struct {
			Type string
		}
		Version struct {
			Type    string
			Pattern string
		}
	}
	Required             []string
	AdditionalProperties bool
}
