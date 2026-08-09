package api

import (
	"encoding/json"
	"groupie-tracker/models"
	"io"
	"net/http"
)

func GetArtists(url string) ([]models.Artists, error) {

	var artists []models.Artists

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, &artists)
	if err != nil {
		return nil, err
	}
	return artists, nil
}
