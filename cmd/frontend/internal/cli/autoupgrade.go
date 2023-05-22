package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/httpapi"
	apirouter "github.com/sourcegraph/sourcegraph/cmd/frontend/internal/httpapi/router"
	"github.com/sourcegraph/sourcegraph/internal/conf"
	"github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/locker"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/multiversion"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/runner"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/store"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/httpserver"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration/migrations"
	"github.com/sourcegraph/sourcegraph/internal/version"
	"github.com/sourcegraph/sourcegraph/internal/version/upgradestore"
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/sourcegraph/lib/output"
)

var buffer strings.Builder // :)

func tryAutoUpgrade(ctx context.Context, logger log.Logger, db database.DB, hook store.RegisterMigratorsUsingConfAndStoreFactoryFunc) (err error) {
	autoupgradeStore := upgradestore.New(db)
	locker := locker.NewWith(db, "autoupgrade")
	_, unlock, err := locker.Lock(ctx, 1, true)
	if err != nil {
		return errors.Wrap(err, "locker.Lock")
	}
	defer func() {
		err = unlock(err)
	}()

	toVersion, _, ok := oobmigration.NewVersionAndPatchFromString(version.Version())
	if !ok {
		return nil
	}
	currentVersionStr, doAutoUpgrade, err := autoupgradeStore.GetAutoUpgrade(ctx)
	if err != nil {
		return errors.Wrap(err, "autoupgradestore.GetAutoUpgrade")
	}
	if !doAutoUpgrade {
		return nil
	}

	currentVersion, _, ok := oobmigration.NewVersionAndPatchFromString(currentVersionStr)
	if !ok {
		return nil
	}

	stopFunc, err := serveConfigurationServer(ctx, logger)
	if err != nil {
		return err
	}
	defer stopFunc()

	if err := runMigration(ctx, currentVersion, toVersion, db, hook); err != nil {
		return err
	}

	return errors.New("MIGRATION SUCCEEDED, RESTARTING")
}

func runMigration(ctx context.Context, from, to oobmigration.Version, db database.DB, hook store.RegisterMigratorsUsingConfAndStoreFactoryFunc) error {
	versionRange, err := oobmigration.UpgradeRange(from, to)
	if err != nil {
		return err
	}

	interrupts, err := oobmigration.ScheduleMigrationInterrupts(from, to)
	if err != nil {
		return err
	}

	plan, err := multiversion.PlanMigration(from, to, versionRange, interrupts)
	if err != nil {
		return err
	}

	registerMigrators := store.ComposeRegisterMigratorsFuncs(
		migrations.RegisterOSSMigratorsUsingConfAndStoreFactory,
		hook,
	)

	return multiversion.RunMigration(ctx,
		db,
		nil, // TODO: nsc do
		plan,
		runner.ApplyPrivilegedMigrations,
		nil,
		true,
		true,
		false,
		true,
		false,
		registerMigrators,
		nil, // only needed for drift
		output.NewOutput(os.Stdout, output.OutputOpts{}), // TODO: nsc do
	)
}

func serveConfigurationServer(ctx context.Context, logger log.Logger) (context.CancelFunc, error) {
	serveMux := http.NewServeMux()
	router := mux.NewRouter().PathPrefix("/.internal").Subrouter()
	middleware := httpapi.JsonMiddleware(&httpapi.ErrorHandler{
		Logger:       logger,
		WriteErrBody: true,
	})
	router.Get(apirouter.Configuration).Handler(middleware(func(w http.ResponseWriter, r *http.Request) error {
		configuration := conf.Unified{
			ServiceConnectionConfig: conftypes.ServiceConnections{
				PostgresDSN:          "lol",
				CodeIntelPostgresDSN: "lol",
				CodeInsightsDSN:      "lol",
			},
		}
		return json.NewEncoder(w).Encode(configuration)
	}))
	serveMux.Handle("/.internal/", router)
	h := http.Handler(serveMux)
	server := &http.Server{
		Handler:      h,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	listener, err := httpserver.NewListener(httpAddrInternal)
	if err != nil {
		return nil, err
	}
	confServer := httpserver.New(listener, server)

	goroutine.Go(func() {
		confServer.Start()
	})

	return confServer.Stop, nil
}

func serveUpgradeUI(ctx context.Context, logger log.Logger) (context.CancelFunc, error) {
	serveMux := http.NewServeMux()
	router := mux.NewRouter().PathPrefix("/.internal").Subrouter()
	middleware := httpapi.JsonMiddleware(&httpapi.ErrorHandler{
		Logger:       logger,
		WriteErrBody: true,
	})
	router.Get(apirouter.Configuration).Handler(middleware(func(w http.ResponseWriter, r *http.Request) error {
		configuration := conf.Unified{
			ServiceConnectionConfig: conftypes.ServiceConnections{
				PostgresDSN:          "lol",
				CodeIntelPostgresDSN: "lol",
				CodeInsightsDSN:      "lol",
			},
		}
		return json.NewEncoder(w).Encode(configuration)
	}))
	serveMux.Handle("/.internal/", router)
	h := http.Handler(serveMux)
	server := &http.Server{
		Handler:      h,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	listener, err := httpserver.NewListener(httpAddrInternal)
	if err != nil {
		return nil, err
	}
	confServer := httpserver.New(listener, server)

	goroutine.Go(func() {
		confServer.Start()
	})

	return confServer.Stop, nil
}
