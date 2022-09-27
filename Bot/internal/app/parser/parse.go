package parser

import "strconv"

// Vacancies stores all vacancies.
type Vacancies struct {
	Items []Vacancy `json:"items"`
}

// Vacancy stores all data about vacancy.
type Vacancy struct {
	Name     string   `json:"name"`
	Area     Area     `json:"area"`
	Salary   Salary   `json:"salary"`
	URL      string   `json:"alternate_url"`
	Employer Employer `json:"employer"`
}

// Area is exported to add ability to parse City string.
type Area struct {
	City string `json:"name"`
}

// Salary is exported to parse salary string.
type Salary struct {
	From     int    `json:"from"`
	To       int    `json:"to"`
	Currency string `json:"currency"`
}

// Employer is exported to add ability to parse employer's name string.
type Employer struct {
	Name string `json:"name"`
}

// ParseSalary parses salary string to form request to HH API.
func ParseSalary(salary Salary) string {
	from := strconv.Itoa(salary.From)
	if from == "0" {
		from = ""
	} else {
		from = "от " + from + " "
	}

	to := strconv.Itoa(salary.To)
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

// ParseExp parses experience string to form request to HH API.
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
