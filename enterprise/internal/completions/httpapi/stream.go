package httpapi

import (
	"context"
	"net/http"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/completions/types"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/redispool"
	streamhttp "github.com/sourcegraph/sourcegraph/internal/search/streaming/http"
	"github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/schema"
)

// NewChatCompletionsStreamHandler is an http handler which streams back completions results.
func NewChatCompletionsStreamHandler(logger log.Logger, db database.DB) http.Handler {
	rl := NewRateLimiter(db, redispool.Store, types.CompletionsFeatureChat)
	return newCompletionsHandler(rl, "stream", func(requestParams types.CompletionRequestParameters, c *schema.Completions) string {
		// No user defined models for now.
		return c.ChatModel
	}, func(ctx context.Context, requestParams types.CompletionRequestParameters, cc types.CompletionsClient, w http.ResponseWriter) {
		eventWriter, err := streamhttp.NewWriter(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Always send a final done event so clients know the stream is shutting down.
		defer func() {
			_ = eventWriter.Event("done", map[string]any{})
		}()

		err = cc.Stream(ctx, types.CompletionsFeatureChat, requestParams, func(event types.CompletionResponse) error { return eventWriter.Event("completion", event) })
		if err != nil {
			trace.Logger(ctx, logger).Error("error while streaming completions", log.Error(err))
			_ = eventWriter.Event("error", map[string]string{"error": err.Error()})
			return
		}
	})
}
