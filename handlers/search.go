package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/poylie/yelp-search/handlers/structs"

	"github.com/Sirupsen/logrus"
	"github.com/poylie/yelp-search/libhttp"
)

// GetSearch handler
func GetSearch(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/html")

	tmpl, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")

	responseObject := GetResult(responseWriter, request)

	responseObject = ProcessSortList(responseObject)
	responseObject = ProcessRating(responseObject)

	tmpl.Execute(responseWriter, responseObject)
}

// GetResult func to get result from API
func GetResult(responseWriter http.ResponseWriter, request *http.Request) structs.Response {
	var responseObject structs.Response

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

	return responseObject
}

// ProcessSortList to identify which sort type was selected
func ProcessSortList(response structs.Response) structs.Response {
	tempMap := map[string]string{
		"":             "Sort By",
		"best_match":   "Relevance",
		"rating":       "Rating",
		"review_count": "Review Count",
	}

	response.SortMap = make(map[string]structs.SortSelection)

	for key, value := range tempMap {
		if key == response.SortBy {
			response.SortMap[key] = structs.SortSelection{value, true}
		} else {
			response.SortMap[key] = structs.SortSelection{value, false}
		}
	}
	return response
}

// ProcessRating process display of rating
func ProcessRating(response structs.Response) structs.Response {
	// This section is for the star rating display
	for elementIndex, element := range response.Business {
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

		response.Business[elementIndex] = element
	}

	return response
}
