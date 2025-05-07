package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const unknown = "uncnown"

func GetGender(name string) string {
	var res struct {
		Gender string `json:"gender"`
	}

	apiURL := os.Getenv("GENDERIZE_API")
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return ""
	}

	return res.Gender
}

func GetAge(name string) int {
	var res struct {
		Age int `json:"age"`
	}

	apiURL := os.Getenv("AGIFY_API")
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0
	}

	return res.Age
}

func GetNationality(name string) string {
	var res struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	apiURL := os.Getenv("NATIONALIZE_API")
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
	if err != nil {
		return unknown
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return unknown
	}

	if len(res.Country) > 0 {
		return res.Country[0].CountryID
	}

	return unknown
}
