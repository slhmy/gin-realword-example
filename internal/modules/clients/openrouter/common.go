package openrouter

import (
	"log/slog"

	"gin-realword-example/internal/modules/core"
	"gin-realword-example/internal/modules/shared"
)

const (
	baseUrl = "https://openrouter.ai"
)

var (
	accessToken string
	model       string
)

func init() {
	accessToken = core.ConfigStore.GetString(shared.ConfigKeyOpenRouterAccessToken)
	model = core.ConfigStore.GetString(shared.ConfigKeyOpenRouterModel)
	if len(accessToken) == 0 {
		slog.Warn("OpenRouter access token is empty")
	}
	if len(model) == 0 {
		slog.Warn("OpenRouter model is empty")
	}
}

type messageObject struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type streamChoiceObject struct {
	Index        int           `json:"index"`
	Delta        messageObject `json:"delta"`
	FinishReason *string       `json:"finish_reason"`
}
