package adapterhttp

type PersonResponse struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

type CreatePersonRequest struct {
	Name       string `json:"name" example:"Dmitriy"`
	Surname    string `json:"surname" example:"Ushakov"`
	Patronymic string `json:"patronymic,omitempty" example:"Vasilevich"`
}
