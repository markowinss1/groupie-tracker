package handlers

import (
	"fmt"
	"groupie-tracker/api"
	"groupie-tracker/models"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

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

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		fmt.Println("Template parsing error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nartist)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

}
