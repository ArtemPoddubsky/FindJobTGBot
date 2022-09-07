package app

import (
	"encoding/json"
	"fmt"
	"io"
	"main/internal/app/parser"
	"main/internal/storage"
	"main/internal/utils"
	"net/http"
	"net/url"
)

const Repeat = "Используйте /repeat чтобы посмотреть следующие 7 вакансий"
const VacPerRequest = 7

func (a *App) StartSearch(id int64, text string) {
	a.Mutex.RLock()
	req := a.Req[id]
	delete(a.Req, id)
	a.Mutex.RUnlock()
	if err := a.Search(id, req, parser.ParseExp(text)); err != nil {
		utils.Error(fmt.Errorf("Search: %v ", err), text)
	}
	if err := a.SendMessage(id, Repeat); err != nil {
		utils.Error(fmt.Errorf("Send /repeat prompt: %v ", err), text)
	}
}

func (a *App) Search(chatID int64, title, exp string) error {
	params := url.Values{}
	params.Add("text", title)
	params.Add("area", "1")
	if exp != "" {
		params.Add("experience", exp)
	}
	resp, err := http.Get("https://api.hh.ru/vacancies?" + params.Encode())
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = a.SearchVacancies(body, chatID, a.db); err != nil {
		return err
	}
	return a.db.AddLast(chatID, title, exp)
}

func (a *App) SearchVacancies(body []byte, chatID int64, db *storage.Postgres) error {
	var s = new(parser.Vacancies)
	if err := json.Unmarshal(body, &s); err != nil {
		return err
	}
	count := 0
	for i := range s.Items {
		if count == VacPerRequest {
			break
		}
		unique, err := db.CheckOriginal(chatID, s.Items[i].Url)
		if err != nil {
			return err
		} else if !unique {
			continue
		}
		err = db.AddRecord(chatID, s.Items[i].Url)
		if err != nil {
			return err
		}
		if err = a.SendVacancy(chatID, s.Items[i], parser.ParseSalary(s.Items[i].Salary)); err != nil {
			return err
		}
		count++
	}
	if count == 0 {
		return a.SendMessage(chatID, "По данному запросу не найдено новых вакансий")
	}
	return nil
}
