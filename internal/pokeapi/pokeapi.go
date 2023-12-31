package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Pertsaa/pokedex-cli/internal/pokecache"
)

type PokeApi struct {
	cache *pokecache.Cache
}

func New(cacheDuration time.Duration) *PokeApi {
	return &PokeApi{
		cache: pokecache.New(cacheDuration),
	}
}

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

func (p *PokeApi) GetLocationAreas(pageURL *string) (*LocationAreaResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pageURL != nil {
		url = *pageURL
	}

	bytes, err := getWithCache(url, p.cache)
	if err != nil {
		return &LocationAreaResponse{}, err
	}

	var lr LocationAreaResponse
	if err := json.Unmarshal(bytes, &lr); err != nil {
		return &LocationAreaResponse{}, err
	}

	return &lr, nil
}

type AreaEncountersResponse struct {
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (p *PokeApi) GetAreaEncounters(area string) (*AreaEncountersResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area

	bytes, err := getWithCache(url, p.cache)
	if err != nil {
		return &AreaEncountersResponse{}, err
	}

	var ae AreaEncountersResponse
	if err := json.Unmarshal(bytes, &ae); err != nil {
		return &AreaEncountersResponse{}, err
	}

	return &ae, nil
}

type PokemonResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func (p *PokeApi) GetPokemon(name string) (*PokemonResponse, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	bytes, err := getWithCache(url, p.cache)
	if err != nil {
		return &PokemonResponse{}, err
	}

	var pokemon PokemonResponse
	if err := json.Unmarshal(bytes, &pokemon); err != nil {
		return &PokemonResponse{}, err
	}

	return &pokemon, nil
}

func getWithCache(url string, cache *pokecache.Cache) ([]byte, error) {
	var bytes []byte

	if cb, ok := cache.Get(url); ok {
		bytes = cb
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return []byte{}, err
		}

		defer resp.Body.Close()

		bytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}

		cache.Add(url, bytes)
	}

	return bytes, nil
}
