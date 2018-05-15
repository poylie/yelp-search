package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/poylie/yelp-search/libhttp"
)

// Response is the entire body of the YELP API Response
type Response struct {
	Business []Business `json:"businesses"`
	Total    int        `json:"total"`
	Keyword  string
	Location string
	SortBy   string
	SortMap  map[string]SortSelection
}

// Sort Selection struct
type SortSelection struct {
	SortDisplay string
	Selected    bool
}

// Business struct for each business
type Business struct {
	Name          string   `json:"name"`
	Alias         string   `json:"alias"`
	Location      Location `json:"location"`
	Price         string   `json:"price"`
	Rating        float64  `json:"rating"`
	RatingDisplay []bool
	RatingHalf    bool
	ReviewCount   int         `json:"review_count"`
	URL           string      `json:"url"`
	ImageURL      string      `json:"image_url"`
	Coordinates   Coordinates `json:"coordinates"`
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
func GetSearch(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")

	request.ParseForm()
	// get value from search
	keyword := request.FormValue("keyword")
	location := request.FormValue("location")
	sortBy := request.FormValue("sortBy")

	client := &http.Client{}
	yelpRequest, _ := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search", nil)

	urlQuery := yelpRequest.URL.Query()
	if keyword != "" {
		urlQuery.Add("term", keyword)
	}

	if location != "" {
		urlQuery.Add("location", location)
	}

	if sortBy != "" {
		urlQuery.Add("sort_by", sortBy)
	}

	yelpRequest.URL.RawQuery = urlQuery.Encode()

	yelpRequest.Header.Set("Authorization", "Bearer VNBosgLfAisHadcgnVbBUUawb--mf61EurggtiBUDcdG-StSyzuwQ39yk0n5OGMIbKxcjFlhaHh8c0YNIrbUmv-UXRaHtYNqcEAkrQ8c-5UMXUerZj85DoRZqU3yWnYx")

	logrus.Infof("Request: %s\n", yelpRequest.URL)
	response, err := client.Do(yelpRequest)

	var responseObject Response

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		libhttp.HandleErrorJson(responseWriter, err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &responseObject)
	}

	responseObject.Location = location
	responseObject.Keyword = keyword
	responseObject.SortBy = sortBy
	tempMap := map[string]string{
		"":             "Sort By",
		"best_match":   "Relevance",
		"rating":       "Rating",
		"review_count": "Review Count",
	}

	responseObject.SortMap = make(map[string]SortSelection)

	for key, value := range tempMap {
		if key == sortBy {
			responseObject.SortMap[key] = SortSelection{value, true}
		} else {
			responseObject.SortMap[key] = SortSelection{value, false}
		}
	}

	// This section is for the star rating display
	for elementIndex, element := range responseObject.Business {
		for idx := 0; idx < int(element.Rating); idx++ {
			if element.RatingDisplay == nil {
				element.RatingDisplay = []bool{true}
			} else {
				element.RatingDisplay = append(element.RatingDisplay, true)
			}

		}
		if element.Rating != math.Trunc(element.Rating) {
			element.RatingHalf = true
		}

		responseObject.Business[elementIndex] = element
	}

	tmpl.Execute(responseWriter, responseObject)
}
