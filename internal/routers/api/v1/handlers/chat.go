package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"gin-realword-example/internal/modules/clients/openrouter"

	"github.com/gin-gonic/gin"
)

// ChatStream
//
//	@Param		prompt	query	string	true	"prompt"
//	@Produce	text/event-stream
//	@Success	200	{string}	string
//	@Router		/chat/stream [get]
func ChatStream(ginCtx *gin.Context) {
	prompt := ginCtx.Query("prompt")
	slog.InfoContext(ginCtx, "chat_stream",
		slog.String("prompt", prompt),
	)
	stream, err := openrouter.ChatCompletionStream(ginCtx.Request.Context(), prompt)
	if err != nil {
		ginCtx.Error(err)
		return
	}

	ginCtx.Header("Content-Type", "text/event-stream")
	ginCtx.Header("Cache-Control", "no-cache")
	ginCtx.Stream(func(w io.Writer) bool {
		select {
		case <-ginCtx.Request.Context().Done():
			w.Write([]byte("event: error\ndata: request canceled\n\n"))
			return false
		case res, ok := <-stream:
			if !ok {
				w.Write([]byte("event: error\ndata: stream closed\n\n"))
				return false
			}
			resStr, err := json.Marshal(res)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("event: error\ndata: %s\n\n", err.Error())))
				return false
			}
			w.Write(fmt.Appendf(nil, "data: %s\n\n", resStr))
			ginCtx.Writer.Flush()
		}
		return true
	})
}
