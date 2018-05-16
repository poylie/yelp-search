package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/poylie/yelp-search/handlers"
	"github.com/poylie/yelp-search/handlers/structs"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGetResult(t *testing.T) {
	httpmock.Activate()

	data := url.Values{}
	data.Set("keyword", "starbucks")
	data.Set("location", "singapore")

	req, err := http.NewRequest("POST", "/search", strings.NewReader(data.Encode()))

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.yelp.com/v3/businesses/search",
		httpmock.NewStringResponder(200, `{
			"businesses": [
				{
					"id": "Zztv1vtQ2AjOSUc_12-2pQ",
					"alias": "alt-pizza-singapore",
					"name": "Alt Pizza",
					"image_url": "https://s3-media1.fl.yelpcdn.com/bphoto/bxoR0HLfZfBbzZ7pwIe-WA/o.jpg",
					"is_closed": false,
					"url": "https://www.yelp.com/biz/alt-pizza-singapore?adjust_creative=4lKQIOYo5xxe4hKyyfbvXg&utm_campaign=yelp_api_v3&utm_medium=api_v3_business_search&utm_source=4lKQIOYo5xxe4hKyyfbvXg",
					"review_count": 27,
					"categories": [
						{
							"alias": "pizza",
							"title": "Pizza"
						},
						{
							"alias": "bars",
							"title": "Bars"
						}
					],
					"rating": 4,
					"coordinates": {
						"latitude": 1.29436062295676,
						"longitude": 103.859161086699
					},
					"transactions": [],
					"price": "$$",
					"location": {
						"address1": "9 Temasek Boulevard",
						"address2": "#01-602",
						"address3": "",
						"city": "Singapore",
						"zip_code": "038989",
						"country": "SG",
						"state": "SG",
						"display_address": [
							"9 Temasek Boulevard",
							"#01-602",
							"Singapore 038989",
							"Singapore"
						]
					},
					"phone": "+6568369207",
					"display_phone": "+65 6836 9207",
					"distance": 2473.088135913626
				}
			],
			"total": 1,
			"region": {
				"center": {
					"longitude": 103.84895324707031,
					"latitude": 1.3141221482579457
				}
			}
		}`))

	hf := handlers.GetResult(recorder, req)
	expectedTotal := 1
	expectedName := "Alt Pizza"

	if hf.Total != expectedTotal {
		t.Errorf("Get Result returned wrong total: got %v want %v",
			hf.Total, expectedTotal)
	}

	if hf.Business[0].Name != expectedName {
		t.Errorf("Get Result returned wrong name: got %v want %v",
			hf.Business[0].Name, expectedName)
	}
}

func TestProcessRatingWholeNumber(t *testing.T) {

	sampleBusiness := structs.Business{
		Name:   "Sample",
		Rating: 4,
	}
	baseResponse := structs.Response{
		Business: []structs.Business{sampleBusiness},
	}

	response := handlers.ProcessRating(baseResponse)
	expectedSize := 4

	if len(response.Business[0].RatingDisplay) != expectedSize {
		t.Errorf("Process Result returned wrong array: got %v want %v",
			len(response.Business[0].RatingDisplay), expectedSize)
	}

}

func TestProcessRatingDecimalNumber(t *testing.T) {

	sampleBusiness := structs.Business{
		Name:   "Sample",
		Rating: 3.5,
	}
	baseResponse := structs.Response{
		Business: []structs.Business{sampleBusiness},
	}

	response := handlers.ProcessRating(baseResponse)
	expectedSize := 3
	expectedHalf := true

	if len(response.Business[0].RatingDisplay) != expectedSize {
		t.Errorf("Process Result returned wrong array: got %v want %v",
			len(response.Business[0].RatingDisplay), expectedSize)
	}

	if response.Business[0].RatingHalf != expectedHalf {
		t.Errorf("Process Result returned wrong value for half: got %v want %v",
			response.Business[0].RatingHalf, expectedHalf)
	}

}

func TestProcessSortList(t *testing.T) {
	baseResponse := structs.Response{
		SortBy: "best_match",
	}

	response := handlers.ProcessSortList(baseResponse)
	expectedSelected := true
	expectedSize := 4

	if len(response.SortMap) != expectedSize {
		t.Errorf("Process Sort returned wrong map size: got %v want %v",
			len(response.SortMap), expectedSize)
	}

	if response.SortMap["best_match"].Selected != expectedSelected {
		t.Errorf("Process Sort returned wrong selected value: got %v want %v",
			response.SortMap["best_match"].Selected, expectedSelected)
	}

}
