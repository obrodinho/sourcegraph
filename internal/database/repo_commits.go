package database

import (
	"context"

	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/log"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/database/batch"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/types"
)

type RepoCommitsStore interface {
	BatchInsertCommitSHAsWithPerforceChangelistID(context.Context, api.RepoID, []types.PerforceChangelist) error
	GetLatestForRepo(ctx context.Context, repoID api.RepoID) (*types.RepoCommit, error)
}

type repoCommitsStore struct {
	*basestore.Store
	logger log.Logger
}

var _ RepoCommitsStore = (*repoCommitsStore)(nil)

func RepoCommitsWith(logger log.Logger, other basestore.ShareableStore) RepoCommitsStore {
	return &repoCommitsStore{
		logger: logger,
		Store:  basestore.NewWithHandle(other.Handle()),
	}
}

func (s *repoCommitsStore) BatchInsertCommitSHAsWithPerforceChangelistID(ctx context.Context, repo_id api.RepoID, commitsMap []types.PerforceChangelist) error {
	return s.WithTransact(ctx, func(tx *basestore.Store) error {

		inserter := batch.NewInserter(ctx, tx.Handle(), "repo_commits", batch.MaxNumPostgresParameters, "repo_id", "commit_sha", "perforce_changelist_id")
		for _, item := range commitsMap {
			if err := inserter.Insert(
				ctx,
				int32(repo_id),
				dbutil.CommitBytea(item.CommitSHA),
				item.ChangelistID,
			); err != nil {
				return err
			}
		}
		return inserter.Flush(ctx)
	})

}

var getLatestForRepoFmtStr = `
SELECT
	id,
	repo_id,
	commit_sha,
	perforce_changelist_id
FROM
	repo_commits
WHERE
	repo_id = %s
ORDER BY
	id DESC
LIMIT 1`

func (s *repoCommitsStore) GetLatestForRepo(ctx context.Context, repoID api.RepoID) (*types.RepoCommit, error) {
	q := sqlf.Sprintf(getLatestForRepoFmtStr, repoID)
	row := s.QueryRow(ctx, q)
	return scanRepoCommitRow(row)
}

func scanRepoCommitRow(scanner dbutil.Scanner) (*types.RepoCommit, error) {
	var r types.RepoCommit
	err := scanner.Scan(
		&r.ID,
		&r.RepoID,
		&r.CommitSHA,
		&r.PerforceChangelistID,
	)
	return &r, err
}
