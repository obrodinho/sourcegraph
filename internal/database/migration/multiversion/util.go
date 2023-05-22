package multiversion

import (
	"context"
	"database/sql"
	"strings"

	"github.com/sourcegraph/sourcegraph/internal/database"
	connections "github.com/sourcegraph/sourcegraph/internal/database/connections/live"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/runner"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/schemas"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/store"
	"github.com/sourcegraph/sourcegraph/internal/database/postgresdsn"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
	"github.com/sourcegraph/sourcegraph/internal/version/upgradestore"
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/sourcegraph/lib/output"
)

const appName = "frontend-autoupgrader"

func NewRunnerWithSchemas(observationCtx *observation.Context, out *output.Output, schemaNames []string, schemas []*schemas.Schema) (*runner.Runner, error) {
	dsns, err := postgresdsn.DSNsBySchema(schemaNames)
	if err != nil {
		return nil, err
	}

	var dsnsStrings []string
	for schema, dsn := range dsns {
		dsnsStrings = append(dsnsStrings, schema+" => "+dsn)
	}

	out.WriteLine(output.Linef(output.EmojiInfo, output.StyleGrey, "Connection DSNs used: %s", strings.Join(dsnsStrings, ", ")))

	storeFactory := func(db *sql.DB, migrationsTable string) connections.Store {
		return connections.NewStoreShim(store.NewWithDB(observationCtx, db, migrationsTable))
	}
	r, err := connections.RunnerFromDSNsWithSchemas(out, observationCtx.Logger, dsns, appName, storeFactory, schemas)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetServiceVersion(ctx context.Context, db database.DB) (oobmigration.Version, int, bool, error) {
	versionStr, ok, err := upgradestore.New(db).GetServiceVersion(ctx)
	if err != nil || !ok {
		return oobmigration.Version{}, 0, ok, err
	}

	version, patch, ok := oobmigration.NewVersionAndPatchFromString(versionStr)
	if !ok {
		return oobmigration.Version{}, 0, ok, errors.Newf("cannot parse version: %q - expected [v]X.Y[.Z]", versionStr)
	}

	return version, patch, true, nil
}
