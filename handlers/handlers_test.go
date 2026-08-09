package handlers

import (
	"testing"

	"groupie-tracker/models"
)

func TestSearchBarFindsMatch(t *testing.T) {
	artists := []models.Artists{
		{Name: "Queen"},
		{Name: "Pink Floyd"},
	}

	result := SearchBar(artists, "queen")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Name != "Queen" {
		t.Errorf("expected Queen, got %s", result[0].Name)
	}
}

func TestSearchBarNoMatch(t *testing.T) {
	artists := []models.Artists{
		{Name: "Queen"},
	}

	result := SearchBar(artists, "zzzzz")
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}

func TestSearchBarCaseInsensitive(t *testing.T) {
	artists := []models.Artists{
		{Name: "Pink Floyd"},
	}

	result := SearchBar(artists, "PINK")
	if len(result) != 1 {
		t.Errorf("expected case-insensitive match, got %d results", len(result))
	}
}
