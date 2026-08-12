package handlers

import (
	"fmt"
	"groupie-tracker/api"
	"groupie-tracker/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type ArtistPageData struct {
	Artist   models.Artists
	Date     models.Dates
	Location models.Locations
}

func SearchBar(artists []models.Artists, s string) []models.Artists {

	var foundArtists []models.Artists
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(s)) {
			foundArtists = append(foundArtists, artist)
		}
	}
	return foundArtists
}

func Mainpage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	artists, err := api.GetArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to load artists list, please try again later", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Template parsing error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		artists = SearchBar(artists, search)
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

}

func Artistpage(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to load artist data", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	var nartist models.Artists
	flag := true
	for _, artist := range artists {
		if artist.Id == id {
			nartist = artist
			flag = false
		}
	}
	if flag {
		http.NotFound(w, r)
		return
	}

	artistDate, err := api.GetDates(nartist.ConcertDates)
	if err != nil {
		http.Error(w, "Error!", http.StatusInternalServerError)
		return
	}
	artistLocation, err := api.GetLocations(nartist.Locations)
	if err != nil {
		http.Error(w, "Error!", http.StatusInternalServerError)
		return
	}

	var artistPage = ArtistPageData{Artist: nartist, Date: artistDate, Location: artistLocation}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		fmt.Println("Template parsing error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, artistPage)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

}
