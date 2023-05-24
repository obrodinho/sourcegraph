package guardrails

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/enterprise"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/guardrails/resolvers"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel"
	"github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

func Init(
	_ context.Context,
	_ *observation.Context,
	_ database.DB,
	_ codeintel.Services,
	_ conftypes.UnifiedWatchable,
	enterpriseServices *enterprise.Services,
) error {
	enterpriseServices.GuardrailsResolver = resolvers.NewGuardrailsResolver()

	return nil
}
