package entity

type Car struct {
	ID     string
	Mark   string   `json:"mark"`   // Марка автомобиля
	Model  string   `json:"model"`  // Модель автомобиля
	RegNum string   `json:"regNum"` // Государственный номер автомобиля
	Year   int      `json:"year"`   // Дата выпуска автомобиля
	Owner  struct { // Владелец авто
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	} `json:"owner"`
}

func NewCar() Car {
	return Car{}
}

type RegNums struct {
	RegNum []string `json:"regNums"`
}
