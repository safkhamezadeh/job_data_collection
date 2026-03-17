package config

import "os"

type ApiKeyConfig struct {
	KeywordKeys    KeywordTranslationConfig
	JobVacancyKeys JobVacancyConfig
}

func LoadApiKeys() ApiKeyConfig {
	return ApiKeyConfig{KeywordKeys: KeywordTranslationConfig{GoogleGeminiApiKey: os.Getenv("GEMINI_API_KEY")}, JobVacancyKeys: JobVacancyConfig{loadAdzunaKeys()}}
}

type KeywordTranslationConfig struct {
	GoogleGeminiApiKey string
}

type JobVacancyConfig struct {
	AdzunaKeys AdzunaKeys
}

type AdzunaKeys struct {
	Adzuna_application_key string
	Adzuna_application_id  string
}

func loadAdzunaKeys() AdzunaKeys {
	return AdzunaKeys{Adzuna_application_key: os.Getenv("ADZUNA_API_APPLICATION_KEY"), Adzuna_application_id: os.Getenv("ADZUNA_API_APPLICATION_ID")}

}

type HttpConfig struct {
	PORT int
}

func LoadHttpConfig() HttpConfig {
	return HttpConfig{PORT: 8080}
}
