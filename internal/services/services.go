package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetGender(name string) string {
	var res struct {
		Gender string `json:"gender"`
	}

	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "" // или логгируй ошибку, или возвращай "unknown"
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

	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
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

	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "unknown"
	}

	if len(res.Country) > 0 {
		return res.Country[0].CountryID
	}

	return "unknown"
}
