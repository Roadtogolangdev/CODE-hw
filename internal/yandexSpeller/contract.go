package yandexSpeller

type Speller interface {
	CheckSpelling(text string) ([]SpellCheckResult, error)
}
