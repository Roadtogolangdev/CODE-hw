package yandexSpeller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SpellCheckResult struct {
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

func CheckSpelling(text string) ([]SpellCheckResult, error) {
	apiUrl := "https://speller.yandex.net/services/spellservice.json/checkText"
	reqUrl := fmt.Sprintf("%s?text=%s", apiUrl, url.QueryEscape(text))

	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []SpellCheckResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
