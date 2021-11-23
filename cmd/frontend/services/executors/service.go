package executors

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/services/executors/store"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/services/executors/store/postgres"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/types"
)

type Executor interface {
	List(ctx context.Context, query string, active bool, offset int, limit int) ([]types.Executor, int, error)
	GetByID(ctx context.Context, gqlID graphql.ID) (types.Executor, bool, error)
	Heartbeat(ctx context.Context, executor types.Executor) error
}

func New(db dbutil.DB) Executor {
	store := postgres.New(db)
	return &executorService{store: store}
}

type executorService struct {
	store store.ExecutorStore
}

func (s *executorService) List(ctx context.Context, query string, active bool, offset int, limit int) ([]types.Executor, int, error) {
	args := store.ExecutorStoreListOptions{
		Query:  query,
		Active: active,
		Offset: offset,
		Limit:  limit,
	}

	return s.store.List(ctx, args)
}

func (s *executorService) GetByID(ctx context.Context, gqlID graphql.ID) (types.Executor, bool, error) {
	id, err := unmarshalExecutorID(gqlID)
	if err != nil {
		return types.Executor{}, false, err
	}

	return s.store.GetByID(ctx, int(id))
}

func (s *executorService) Heartbeat(ctx context.Context, executor types.Executor) error {
	return s.store.UpsertHeartbeat(ctx, executor)
}

func unmarshalExecutorID(id graphql.ID) (executorID int64, err error) {
	err = relay.UnmarshalSpec(id, &executorID)
	return
}
