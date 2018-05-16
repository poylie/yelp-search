package structs

// Response is the entire body of the YELP API Response
type Response struct {
	Business []Business `json:"businesses"`
	Total    int        `json:"total"`
	Keyword  string
	Location string
	SortBy   string
	SortMap  map[string]SortSelection
}

// SortSelection struct
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
