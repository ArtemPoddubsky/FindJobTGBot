package app

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"main/internal/app/parser"
	"main/internal/utils"
	"net/http"
	"net/url"
)

const Repeat = "Используйте /repeat чтобы посмотреть следующие 7 вакансий"
const VacPerRequest = 7

func (a *App) StartSearch(id int64, text string) { // Search
	a.Mutex.RLock()
	req := a.Req[id]
	delete(a.Req, id)
	a.Mutex.RUnlock()
	if err := a.Search(id, req, parser.ParseExp(text)); err != nil {
		utils.FieldError("Search:", err, text)
	}
	if err := a.SendMessage(id, Repeat); err != nil {
		log.Errorln("Sending Repeat prompt", err)
	}
}

func (a *App) Search(chatID int64, title, exp string) error {
	s, err := GetVacancies(title, exp)
	if err != nil {
		return err
	}
	count := 0
	for i := range s.Items {
		if count == VacPerRequest {
			break
		}
		unique, err2 := a.db.CheckOriginal(chatID, s.Items[i].Url)
		if err2 != nil {
			return err2
		} else if !unique {
			continue
		}
		err2 = a.db.AddRecord(chatID, s.Items[i].Url)
		if err2 != nil {
			return err2
		}
		if err2 = a.SendVacancy(chatID, s.Items[i], parser.ParseSalary(s.Items[i].Salary)); err != nil {
			return err
		}
		count++
	}
	if count == 0 {
		return a.SendMessage(chatID, "По данному запросу не найдено новых вакансий")
	}
	return a.db.AddLast(chatID, title, exp)
}

func GetVacancies(title, exp string) (parser.Vacancies, error) {
	params := url.Values{}
	params.Add("text", title)
	params.Add("area", "1")
	if exp != "" {
		params.Add("experience", exp)
	}
	var s parser.Vacancies
	resp, err := http.Get("https://api.hh.ru/vacancies?" + params.Encode())
	if err != nil {
		return s, err
	}
	return s, json.NewDecoder(resp.Body).Decode(&s)
}
