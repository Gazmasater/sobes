package serv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ExternalService interface {
	GetAge(name string) int
	GetGender(name string) string
	GetNationality(name string) string
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

func (es *ExternalServiceImpl) GetAge(name string) int {

	fmt.Printf("GetAge NAME=%s\n", name)
	fmt.Printf("GetAge API=%s\n", es.AgifyAPI)

	url := fmt.Sprintf("%s?name=%s", es.AgifyAPI, name)
	fmt.Printf("GetAge URL=%s\n", url)

	resp, err := http.Get(url)
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
	fmt.Printf("GetAge AGE=%d\n", result.Age)

	return result.Age
}

func (es *ExternalServiceImpl) GetGender(name string) string {
	url := fmt.Sprintf("%s?name=%s", es.GenderizeAPI, name)
	resp, err := http.Get(url)
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

func (es *ExternalServiceImpl) GetNationality(name string) string {
	url := fmt.Sprintf("%s?name=%s", es.NationalizeAPI, name)
	fmt.Printf("GetNationality URL=%s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Nationality string `json:"country_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Nationality
}
