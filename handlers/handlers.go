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

type MainpageInfo struct {
	Artists     []models.Artists
	SearchQuery string
}

type ArtistPageData struct {
	Artist      models.Artists
	Date        models.Dates
	Location    models.Locations
	SearchQuery string
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
		renderError(w, http.StatusNotFound, "templates/404.html")
		return
	}

	artists, err := api.GetArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Template parsing error:", err)
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		artists = SearchBar(artists, search)
	}

	var MainpageData = MainpageInfo{Artists: artists, SearchQuery: search}
	err = tmpl.Execute(w, MainpageData)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}

}

func Artistpage(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}

	searchQuery := r.URL.Query().Get("searchQuery")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		renderError(w, http.StatusBadRequest, "templates/400.html")
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
		renderError(w, http.StatusNotFound, "templates/404.html")
		return
	}

	artistDate, err := api.GetDates(nartist.ConcertDates)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}
	artistLocation, err := api.GetLocations(nartist.Locations)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}

	var artistPage = ArtistPageData{Artist: nartist, Date: artistDate, Location: artistLocation, SearchQuery: searchQuery}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		fmt.Println("Template parsing error:", err)
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}
	err = tmpl.Execute(w, artistPage)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "templates/500.html")
		return
	}

}

func renderError(w http.ResponseWriter, statusCode int, templatePath string) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Template parsing error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}
