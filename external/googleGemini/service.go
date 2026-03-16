package googlegemini

import (
	"context"
	"fmt"
	"job_vacancies/internal/keywordextractor"
	"log"

	"google.golang.org/genai"
)

const MODEL string = "gemini-2.5-flash-lite"

const PROMPT string = `Return 5 job titles and 5 keywords.
Description: %s.
Format:
titles: t1, t2, t3, t4, t5
keywords: w1, w2, w3, w4, w5`

type geminiClient struct {
	Client *genai.Client
}

func NewGeminiClient(client *genai.Client) keywordextractor.KeywordsExtractor {
	return &geminiClient{Client: client}
}

func (g *geminiClient) Translate(ctx context.Context, inputstr string) (*keywordextractor.KeyWordFormat, error) {
	fullPrompt := CompleteString(PROMPT, inputstr)
	res, err := g.Client.Models.GenerateContent(ctx, MODEL, genai.Text(fullPrompt), nil)
	if err != nil {
		log.Printf("GeminiClient Translate err: %s", err)
		return nil, keywordextractor.ExternalErr
	}

	fmt := keywordextractor.StoKeyWordFormat(res.Text())

	return &fmt, nil
}

func CompleteString(incomplete string, toAdd string) string {

	return fmt.Sprintf(incomplete, toAdd)
}
