package yandexSpeller

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SpellCheckResult struct {
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

type speller struct {
	client *http.Client
}

func NewSpeller() Speller {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &speller{client: &http.Client{Transport: tr}}
}

func (s *speller) CheckSpelling(text string) ([]SpellCheckResult, error) {
	apiUrl := "https://speller.yandex.net/services/spellservice.json/checkText"
	reqUrl := fmt.Sprintf("%s?text=%s", apiUrl, url.QueryEscape(text))

	resp, err := s.client.Get(reqUrl)

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
