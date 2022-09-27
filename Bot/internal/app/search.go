package app

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"main/internal/app/parser"
	"main/internal/utils"
	"net/http"
	"net/url"
	"time"
)

const (
	repeat        = "Используйте /repeat чтобы посмотреть следующие 7 вакансий"
	vacPerRequest = 7
	reqTimeout    = 2 * time.Second
)

// StartSearch forms request by accessing in-memory storage,
// calls HH API with formed request, sends unique vacancies if found
// and sends prompt for /repeat.
func (a *App) StartSearch(id int64, text string) {
	a.Mutex.RLock()
	req := a.Req[id]
	delete(a.Req, id)
	a.Mutex.RUnlock()

	if err := a.search(id, req, parser.ParseExp(text)); err != nil {
		utils.FieldError("search:", err, text)
	}

	if err := a.sendMessage(id, repeat); err != nil {
		log.Errorln("sendMessage", err)
	}
}

func (a *App) search(chatID int64, title, exp string) error {
	vacancies, errJSON := getVacancies(title, exp)

	if errJSON != nil {
		return fmt.Errorf("getVacancies: %w", errJSON)
	}

	count := 0
	for i := range vacancies.Items {
		if count == vacPerRequest {
			break
		}

		unique, err := a.db.CheckOriginal(chatID, vacancies.Items[i].URL)
		if err != nil {
			return fmt.Errorf("CheckOriginal: %w", err)
		} else if !unique {
			continue
		}

		if err := a.db.AddRecord(chatID, vacancies.Items[i].URL); err != nil {
			return fmt.Errorf("AddRecord: %w", err)
		}

		if err := a.sendVacancy(chatID, &vacancies.Items[i], parser.ParseSalary(vacancies.Items[i].Salary)); err != nil {
			return fmt.Errorf("sendVacancy: %w", err)
		}
		count++
	}

	if count == 0 {
		if err := a.sendMessage(chatID, "По данному запросу не найдено новых вакансий"); err != nil {
			return fmt.Errorf("sendMessage: %w", err)
		}
	}

	if err := a.db.AddLast(chatID, title, exp); err != nil {
		return fmt.Errorf("AddLast: %w", err)
	}

	return nil
}

func getVacancies(title, exp string) (*parser.Vacancies, error) {
	params := url.Values{}
	params.Add("text", title)
	params.Add("area", "1")

	if exp != "" {
		params.Add("experience", exp)
	}

	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, "https://api.hh.ru/vacancies?"+params.Encode(), http.NoBody)

	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", err)
	}
	defer resp.Body.Close()

	var s parser.Vacancies
	err = json.NewDecoder(resp.Body).Decode(&s)

	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &s, nil
}
