package config

import "os"

type Config struct {
	APIKeys APIKeyConfig
	HTTP    HTTPConfig
}

type APIKeyConfig struct {
	KeywordKeys    KeywordTranslationConfig
	JobVacancyKeys JobVacancyConfig
}

type KeywordTranslationConfig struct {
	GoogleGeminiAPIKey string
}

type JobVacancyConfig struct {
	AdzunaKeys AdzunaKeys
}

type AdzunaKeys struct {
	ApplicationKey string
	ApplicationID  string
}

type HTTPConfig struct {
	Port string
}

func Load() Config {
	return Config{
		APIKeys: loadAPIKeys(),
		HTTP:    loadHTTPConfig(),
	}
}

func loadAPIKeys() APIKeyConfig {
	return APIKeyConfig{
		KeywordKeys: KeywordTranslationConfig{
			GoogleGeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		},
		JobVacancyKeys: JobVacancyConfig{
			AdzunaKeys: LoadAdzunaKeys(),
		},
	}
}

func LoadAdzunaKeys() AdzunaKeys {
	return AdzunaKeys{
		ApplicationKey: os.Getenv("ADZUNA_API_APPLICATION_KEY"),
		ApplicationID:  os.Getenv("ADZUNA_API_APPLICATION_ID"),
	}
}

func loadHTTPConfig() HTTPConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return HTTPConfig{Port: port}
}
