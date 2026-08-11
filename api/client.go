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
func GetLocations(url string) (models.Locations, error) {
	var location models.Locations
	resp, err := http.Get(url)
	if err != nil {
		return location, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return location, err
	}

	err = json.Unmarshal(result, &location)
	if err != nil {
		return location, err
	}
	return location, nil
}

func GetDates(url string) (models.Dates, error) {
	var date models.Dates
	resp, err := http.Get(url)
	if err != nil {
		return date, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return date, err
	}

	err = json.Unmarshal(result, &date)
	if err != nil {
		return date, err
	}
	return date, nil
}
