package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	requestTimeout = 5 * time.Second
	unknown        = "unknown"
)

func GetGender(name string) string {
	var res struct {
		Gender string `json:"gender"`
	}

	apiURL := os.Getenv("GENDERIZE_API")
	err := fetchJSON(fmt.Sprintf("%s?name=%s", apiURL, name), &res)
	if err != nil {
		return ""
	}

	return res.Gender
}

func GetAge(name string) int {
	var res struct {
		Age int `json:"age"`
	}

	apiURL := os.Getenv("AGIFY_API")
	err := fetchJSON(fmt.Sprintf("%s?name=%s", apiURL, name), &res)
	if err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", apiURL, name), nil)
	if err != nil {
		return unknown
	}

	resp, err := http.DefaultClient.Do(req)
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

func fetchJSON[T any](url string, target *T) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
