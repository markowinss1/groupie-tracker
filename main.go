package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ApiLinks struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relations string `json:"relations"`
}

type Artists struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	Id             int      `json:"id"`
	DatesLocations []string `json:"dateslocations"`
}

func GetArtists(url string) ([]Artists, error) {

	var artists []Artists

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

func SearchBar(artists *[]Artists, s string) []Artists {

	var foundArtists []Artists
	for _, artist := range *artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(s)) {
			foundArtists = append(foundArtists, artist)
		}
	}
	return foundArtists
}

func Mainpage(w http.ResponseWriter, r *http.Request) {

	artists, err := GetArtists("https://groupietrackers.herokuapp.com/api/artists")
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
		artists = SearchBar(&artists, search)
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

}

func Artistpage(w http.ResponseWriter, r *http.Request) {
	artists, err := GetArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to load artist data", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	var nartist Artists
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

func main() {

	http.HandleFunc("/", Mainpage)
	http.HandleFunc("/artist", Artistpage)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
