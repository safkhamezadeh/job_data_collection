package keywordextractor

import (
	"context"
)

type KeywordsExtractor interface {
	Translate(ctx context.Context, inputstr string) (*KeyWordFormat, error)
}
