package auth

import (
	"net/http"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/enterprise/cmd/llm-proxy/internal/actor"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/llm-proxy/internal/events"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/llm-proxy/internal/response"
	llmproxy "github.com/sourcegraph/sourcegraph/enterprise/internal/llm-proxy"
	"github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type Authenticator struct {
	Logger      log.Logger
	EventLogger events.Logger
	Sources     actor.Sources
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := trace.Logger(r.Context(), a.Logger)
		token, err := ExtractBearer(r.Header)
		if err != nil {
			response.JSONError(logger, w, http.StatusBadRequest, err)
			return
		}

		act, err := a.Sources.Get(r.Context(), token)
		if err != nil {
			response.JSONError(logger, w, http.StatusUnauthorized, err)

			err := a.EventLogger.LogEvent(
				r.Context(),
				events.Event{
					Name:       llmproxy.EventNameUnauthorized,
					Source:     "anonymous",
					Identifier: "anonymous",
				},
			)
			if err != nil {
				logger.Error("failed to log event", log.Error(err))
			}
			return
		}

		if !act.AccessEnabled {
			response.JSONError(
				logger,
				w,
				http.StatusForbidden,
				errors.New("LLM proxy access not enabled"),
			)

			err := a.EventLogger.LogEvent(
				r.Context(),
				events.Event{
					Name:       llmproxy.EventNameAccessDenied,
					Source:     act.Source.Name(),
					Identifier: act.ID,
				},
			)
			if err != nil {
				logger.Error("failed to log event", log.Error(err))
			}
			return
		}

		r = r.WithContext(actor.WithActor(r.Context(), act))
		// Continue with the chain.
		next.ServeHTTP(w, r)
	})
}
