package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Pertsaa/pokedex-cli/internal/pokecache"
)

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokeApi struct {
	cache *pokecache.Cache
}

func New(cacheDuration time.Duration) *PokeApi {
	return &PokeApi{
		cache: pokecache.New(cacheDuration),
	}
}

func (p *PokeApi) GetLocationAreas(pageURL *string) (*LocationAreaResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pageURL != nil {
		url = *pageURL
	}

	var bytes []byte

	if cb, ok := p.cache.Get(url); ok {
		bytes = cb
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return &LocationAreaResponse{}, err
		}

		defer resp.Body.Close()

		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return &LocationAreaResponse{}, err
		}

		p.cache.Add(url, bytes)
	}

	var lr LocationAreaResponse
	err := json.Unmarshal(bytes, &lr)
	if err != nil {
		return &LocationAreaResponse{}, err
	}

	return &lr, nil
}
