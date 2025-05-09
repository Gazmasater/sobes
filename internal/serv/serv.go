package serv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ExternalService interface {
	GetAge(ctx context.Context, name string) int
	GetGender(ctx context.Context, name string) string
	GetNationality(ctx context.Context, name string) string
}

type ExternalServiceImpl struct {
	AgifyAPI       string
	GenderizeAPI   string
	NationalizeAPI string
}

func NewExternalService() *ExternalServiceImpl {
	return &ExternalServiceImpl{
		AgifyAPI:       os.Getenv("AGIFY_API"),
		GenderizeAPI:   os.Getenv("GENDERIZE_API"),
		NationalizeAPI: os.Getenv("NATIONALIZE_API"),
	}
}

func (es *ExternalServiceImpl) GetAge(ctx context.Context, name string) int {
	url := fmt.Sprintf("%s?name=%s", es.AgifyAPI, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var result struct {
		Age int `json:"age"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0
	}

	return result.Age
}

func (es *ExternalServiceImpl) GetGender(ctx context.Context, name string) string {
	url := fmt.Sprintf("%s?name=%s", es.GenderizeAPI, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Gender
}
func (es *ExternalServiceImpl) GetNationality(ctx context.Context, name string) string {
	url := fmt.Sprintf("%s?name=%s", es.NationalizeAPI, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	if len(result.Country) > 0 {
		return result.Country[0].CountryID
	}

	return ""
}
