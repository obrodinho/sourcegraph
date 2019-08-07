package diagnostics

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/diagnostics"
)

// 🚨 SECURITY: TODO!(sqs): there needs to be security checks everywhere here! there are none

// gqlThreadDiagnosticEdge implements the GraphQL type ThreadDiagnosticEdge.
type gqlThreadDiagnosticEdge struct {
	db *dbThreadDiagnostic
}

// threadDiagnosticEdgeByID looks up and returns the ThreadDiagnosticEdge with the given GraphQL
// ID. If no such ThreadDiagnosticEdge exists, it returns a non-nil error.
func threadDiagnosticEdgeByID(ctx context.Context, id graphql.ID) (*gqlThreadDiagnosticEdge, error) {
	dbID, err := graphqlbackend.UnmarshalThreadDiagnosticEdgeID(id)
	if err != nil {
		return nil, err
	}
	return threadDiagnosticEdgeByDBID(ctx, dbID)
}

func (GraphQLResolver) ThreadDiagnosticEdgeByID(ctx context.Context, id graphql.ID) (graphqlbackend.ThreadDiagnosticEdge, error) {
	return threadDiagnosticEdgeByID(ctx, id)
}

// threadDiagnosticEdgeByDBID looks up and returns the ThreadDiagnosticEdge with the given database ID. If
// no such ThreadDiagnosticEdge exists, it returns a non-nil error.
func threadDiagnosticEdgeByDBID(ctx context.Context, dbID int64) (*gqlThreadDiagnosticEdge, error) {
	v, err := dbThreadsDiagnostics{}.GetByID(ctx, dbID)
	if err != nil {
		return nil, err
	}
	return &gqlThreadDiagnosticEdge{db: v}, nil
}

func (v *gqlThreadDiagnosticEdge) ID() graphql.ID {
	return graphqlbackend.MarshalThreadDiagnosticEdgeID(v.db.ID)
}

func (v *gqlThreadDiagnosticEdge) Thread(ctx context.Context) (graphqlbackend.Thread, error) {
	return graphqlbackend.ThreadByID(ctx, graphqlbackend.MarshalThreadID(v.db.ThreadID))
}

func (v *gqlThreadDiagnosticEdge) Diagnostic() (graphqlbackend.Diagnostic, error) {
	return diagnostics.NewGQLDiagnostic(v.db.Type, v.db.Data), nil
}

func (v *gqlThreadDiagnosticEdge) ViewerCanUpdate(ctx context.Context) (bool, error) {
	thread, err := v.Thread(ctx)
	if err != nil {
		return false, err
	}
	return thread.ViewerCanUpdate(ctx)
}
