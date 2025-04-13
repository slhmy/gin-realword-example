package openrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type ChatCompletionStreamResponse struct {
	Content string `json:"content"`
	Error   error  `json:"error"`
	Done    bool   `json:"done"`
}

type chatCompletionRequest struct {
	Model    string          `json:"model"`
	Messages []messageObject `json:"messages"`
	Stream   bool            `json:"stream"`
}

type chatCompletionStreamChunk struct {
	Choices []streamChoiceObject `json:"choices"`
}

func ChatCompletionStream(
	ctx context.Context, prompt string,
) (<-chan ChatCompletionStreamResponse, error) {
	url, err := url.JoinPath(baseUrl, "api/v1/chat/completions")
	if err != nil {
		return nil, err
	}
	payloadObject := chatCompletionRequest{
		Model: model,
		Messages: []messageObject{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: true,
	}
	payload, err := json.Marshal(payloadObject)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	responseChannel := make(chan ChatCompletionStreamResponse)
	go func() {
		defer close(responseChannel)
		defer func() {
			err := res.Body.Close()
			if err != nil {
				slog.Error("Failed to close response body", "error", err)
			}
		}()

		var buffer bytes.Buffer
		reader := bufio.NewReader(res.Body)

		for {
			select {
			case <-ctx.Done():
				responseChannel <- ChatCompletionStreamResponse{
					Error: ctx.Err(),
					Done:  true,
				}
				return
			default:
				lineBytes, err := reader.ReadBytes('\n')
				if err != nil {
					if err.Error() == "EOF" {
						responseChannel <- ChatCompletionStreamResponse{
							Content: buffer.String(),
							Done:    true,
						}
						return
					}
					responseChannel <- ChatCompletionStreamResponse{
						Error: err,
						Done:  true,
					}
					return
				}

				line := strings.TrimSpace(string(lineBytes))
				if !strings.HasPrefix(line, "data: ") {
					continue
				}

				line = strings.TrimPrefix(line, "data: ")
				if line == "[DONE]" {
					responseChannel <- ChatCompletionStreamResponse{
						Done: true,
					}
					return
				}
				var chunk chatCompletionStreamChunk
				err = json.NewDecoder(strings.NewReader(line)).Decode(&chunk)
				if err != nil {
					slog.Error("Failed to decode JSON", "error", err, "line", line)
					continue
				}
				if len(chunk.Choices) == 0 {
					continue
				}

				done := false
				if chunk.Choices[0].FinishReason != nil {
					done = true
				}
				content := chunk.Choices[0].Delta.Content
				buffer.WriteString(content)
				responseChannel <- ChatCompletionStreamResponse{
					Content: content,
					Done:    done,
				}
			}
		}
	}()

	return responseChannel, nil
}
