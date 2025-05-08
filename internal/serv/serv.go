package serv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ExternalService интерфейс для внешних сервисов
type ExternalService interface {
	GetAge(name string) int
	GetGender(name string) string
	GetNationality(name string) string
}

// ExternalServiceImpl структура, которая реализует интерфейс ExternalService
type ExternalServiceImpl struct {
	AgifyAPI       string
	GenderizeAPI   string
	NationalizeAPI string
}

// NewExternalService создает новый экземпляр ExternalService с API URL
func NewExternalService() *ExternalServiceImpl {
	return &ExternalServiceImpl{
		AgifyAPI:       os.Getenv("AGIFY_API"),
		GenderizeAPI:   os.Getenv("GENDERIZE_API"),
		NationalizeAPI: os.Getenv("NATIONALIZE_API"),
	}
}

// GetAge получает возраст по имени через API Agify
func (es *ExternalServiceImpl) GetAge(name string) int {
	url := fmt.Sprintf("%s?name=%s", es.AgifyAPI, name)
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

	return result.Age
}

// GetGender получает пол по имени через API Genderize
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

// GetNationality получает национальность по имени через API Nationalize
func (es *ExternalServiceImpl) GetNationality(name string) string {
	url := fmt.Sprintf("%s?name=%s", es.NationalizeAPI, name)
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
