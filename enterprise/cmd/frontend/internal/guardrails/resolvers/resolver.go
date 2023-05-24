package resolvers

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend/graphqlutil"
)

var _ graphqlbackend.GuardrailsResolver = &guardrailsResolver{}

// TODO(keegancsmith) stub implementation
type guardrailsResolver struct{}

func NewGuardrailsResolver() graphqlbackend.GuardrailsResolver {
	return &guardrailsResolver{}
}

func (c *guardrailsResolver) SnippetAttribution(ctx context.Context, args *graphqlbackend.SnippetAttributionArgs) (*graphqlbackend.SnippetAttributionConnectionResolver, error) {
	return &graphqlbackend.SnippetAttributionConnectionResolver{
		TotalCount: 1,
		LimitHit:   false,
		PageInfo:   graphqlutil.HasNextPage(false),
		Nodes: []*graphqlbackend.SnippetAttributionResolver{{
			RepositoryName: "horse/graph",
		}},
	}, nil
}
