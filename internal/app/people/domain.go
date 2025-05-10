package people

type Person struct {
	ID          uint
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

type Filter struct {
	Gender      string
	Nationality string
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	SortBy      string
	Order       string
	Limit       int
	Offset      int
}
