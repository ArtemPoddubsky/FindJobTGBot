package parser

import "strconv"

type Vacancies struct {
	Items []Vacancy `json:"items"`
}
type Vacancy struct {
	Name     string   `json:"name"`
	Area     Area     `json:"area"`
	Salary   Salary   `json:"salary"`
	Url      string   `json:"alternate_url"`
	Employer Employer `json:"employer"`
}
type Area struct {
	City string `json:"name"`
}
type Salary struct {
	From     int    `json:"from"`
	To       int    `json:"to"`
	Currency string `json:"currency"`
}
type Employer struct {
	Name string `json:"name"`
}

func ParseSalary(salary Salary) string {
	var from, to string
	from = strconv.Itoa(salary.From)
	to = strconv.Itoa(salary.To)
	if from == "0" {
		from = ""
	} else {
		from = "от " + from + " "
	}
	if to == "0" {
		to = ""
	} else {
		to = "до " + to + " "
	}
	if from == "" && to == "" {
		return "Зарплата: Не указана"
	}
	return "Зарплата: " + from + to + salary.Currency
}

func ParseExp(req string) string {
	switch req {
	case "Нет опыта":
		return "noExperience"
	case "От 1 до 3 лет":
		return "between1And3"
	case "От 3 до 6 лет":
		return "between3And6"
	case "Более 6 лет":
		return "moreThan6"
	}
	return ""
}
