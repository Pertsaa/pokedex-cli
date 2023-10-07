package pokeapi

import (
	"encoding/json"
	"net/http"
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

type PokeApi struct{}

func New() *PokeApi {
	return &PokeApi{}
}

func (pa *PokeApi) GetLocationAreas(pageURL *string) (*LocationAreaResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pageURL != nil {
		url = *pageURL
	}

	resp, err := http.Get(url)
	if err != nil {
		return &LocationAreaResponse{}, err
	}

	defer resp.Body.Close()

	var lr LocationAreaResponse

	err = json.NewDecoder(resp.Body).Decode(&lr)
	if err != nil {
		return &LocationAreaResponse{}, err
	}

	return &lr, nil
}
