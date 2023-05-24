package graphqlbackend

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend/graphqlutil"
)

type GuardrailsResolver interface {
	SnippetAttribution(ctx context.Context, args *SnippetAttributionArgs) (*SnippetAttributionConnectionResolver, error)
}

type SnippetAttributionArgs struct {
	Snippet string
	First   *int
	After   *int
}

type SnippetAttributionConnectionResolver struct {
	TotalCount int
	LimitHit   bool
	PageInfo   *graphqlutil.PageInfo
	Nodes      []*SnippetAttributionResolver
}

type SnippetAttributionResolver struct {
	RepositoryName string
}
