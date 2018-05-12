package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/poylie/yelp-search/libhttp"
)

// Response is the entire body of the YELP API Response
type Response struct {
	Business []Business `json:"businesses"`
	Total    int        `json:"total"`
	Keyword  string
	Location string
}

// Business struct for each business
type Business struct {
	Name        string      `json:"name"`
	Alias       string      `json:"alias"`
	Location    Location    `json:"location"`
	Price       string      `json:"price"`
	Rating      float64     `json:"rating"`
	URL         string      `json:"url"`
	ImageURL    string      `json:"image_url"`
	Coordinates Coordinates `json:"coordinates"`
}

// Coordinates struct
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Location holds address of business
type Location struct {
	Address1       string   `json:"address1"`
	Address2       string   `json:"address2"`
	Address3       string   `json:"address3"`
	City           string   `json:"city"`
	ZipCode        string   `json:"zip_code"`
	Country        string   `json:"country"`
	State          string   `json:"state"`
	DisplayAddress []string `json:"display_address"`
}

// GetSearch handler
func GetSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")

	r.ParseForm()
	// logic part of log in
	keyword := r.FormValue("keyword")
	location := r.FormValue("location")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search", nil)

	q := req.URL.Query()
	q.Add("term", keyword)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer VNBosgLfAisHadcgnVbBUUawb--mf61EurggtiBUDcdG-StSyzuwQ39yk0n5OGMIbKxcjFlhaHh8c0YNIrbUmv-UXRaHtYNqcEAkrQ8c-5UMXUerZj85DoRZqU3yWnYx")

	fmt.Printf("Request: %s\n", req.URL)
	response, err := client.Do(req)

	var responseObject Response

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//fmt.Printf(string(data))

		json.Unmarshal(data, &responseObject)
	}

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	responseObject.Location = location
	responseObject.Keyword = keyword

	tmpl.Execute(w, responseObject)
}
