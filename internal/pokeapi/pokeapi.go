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

func GetLocationAreas(url string) (*LocationAreaResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return &LocationAreaResponse{}, err
	}

	defer resp.Body.Close()

	var locationAreaResponse LocationAreaResponse

	err = json.NewDecoder(resp.Body).Decode(&locationAreaResponse)
	if err != nil {
		return &LocationAreaResponse{}, err
	}

	return &locationAreaResponse, nil
}
