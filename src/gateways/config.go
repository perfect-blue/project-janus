package gateways

import (
	"encoding/json"
	"io"
	"os"
)

// TODO: what is this
type Config struct {
	Version   string        `json:"version"`
	Name      string        `json:"name"`
	Port      int           `json:"port"`
	CacheTTL  string        `json:"cache_ttl"`
	Timeout   string        `json:"timeout"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Endpoint     string            `json:"endpoint"`
	Method       HTTPMethod        `json:"method"`
	Backend      []Backend         `json:"backend"`
	Query        []string          `json:"query"`
	QueryMapping map[string]string `json:"query_mapping"`
}

type Backend struct {
	URLPattern string   `json:"url_pattern"`
	Host       []string `json:"host"`
	Port       int      `json:"port"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}