package internal

import "net/http"

type (
	TerrainService struct {
		Server *http.Server
	}
	TerrainServiceOptions func(service *TerrainService)
)

func NewTerrainService(opts ...TerrainServiceOptions) (*TerrainService, error) {
	// default
	port := ":8080"
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	ts := &TerrainService{
		Server: s,
	}

	for _, opt := range opts {
		opt(ts)
	}

	return &TerrainService{
		Server: s,
	}, nil
}

func ConfigWithPort(port string) TerrainServiceOptions {
	return func(s *TerrainService) {
		s.Server.Addr = port
	}
}

func ConfigWithUrlPrefix(prefix string) TerrainServiceOptions {
	return func(s *TerrainService) {
		s.Server.Addr = prefix
	}
}
