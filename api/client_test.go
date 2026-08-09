package api

import "testing"

func TestGetArtists(t *testing.T) {
	artists, err := GetArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(artists) == 0 {
		t.Fatal("expected non-empty artists list")
	}
	if artists[0].Name == "" {
		t.Error("expected first artist to have a name")
	}
}

func TestGetArtistsInvalidURL(t *testing.T) {
	_, err := GetArtists("http://not-a-real-domain-xyz123.com")
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
