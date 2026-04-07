package googlegemini_test

import (
	"context"
	googlegemini "job_vacancies/external/googleGemini"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func TestTranslate_Integration(t *testing.T) {
	if os.Getenv("DEV") != "true" {
		t.Skip("skipping integration test (DEV != true)")
	}

	err := godotenv.Load("./.env.test")
	if err != nil {
		t.Fatalf("error loading .env: %v", err)
	}

	ctx := context.Background()

	prompt := `Return 5 job titles and 5 keywords.
Description: i want to find jobs in a microbiology lab.
Format:
titles: t1, t2, t3, t4, t5
keywords: w1, w2, w3, w4, w5`

	genclient, err := genai.NewClient(ctx, nil)
	if err != nil {
		t.Fatalf("failed to create genai client: %v", err)
	}
	client := googlegemini.NewGeminiClient(genclient)

	res, err := client.Translate(ctx, prompt)
	if err != nil {
		t.Fatalf("translate error: %v", err)
	}
	if len(res.JobTitles) != 5 {
		t.Errorf("expected 5 jobs, got %d", len(res.JobTitles))
	}
	if len(res.Keywords) != 5 {
		t.Errorf("expected 5 keywords, got %d", len(res.Keywords))
	}
}
