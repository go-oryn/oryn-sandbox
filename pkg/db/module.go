package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/XSAM/otelsql"
	"github.com/go-oryn/oryn-sandbox/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "db"

var Module = fx.Module(
	ModuleName,
	// dependencies
	fx.Provide(
		ProvideDB,
		ProvideMigrator,
		ProvideSeeder,
	),
)

type ProvideDBParams struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Shutdown       fx.Shutdowner
	Config         *config.Config
	Logger         *slog.Logger
	Propagator     propagation.TextMapPropagator
	MeterProvider  metric.MeterProvider
	TracerProvider trace.TracerProvider
}

func ProvideDB(params ProvideDBParams) (*sql.DB, error) {
	driver := params.Config.GetString("db.driver")
	if driver == "" {
		return nil, errors.New("db driver is not configured")
	}

	dsn := params.Config.GetString("db.dsn")
	if dsn == "" {
		return nil, errors.New("db dsn is not configured")
	}

	attrs := append(
		otelsql.AttributesFromDSN(params.Config.GetString("db.dsn")),
		semconv.DBSystemNameMySQL,
	)

	db, err := otelsql.Open(
		driver,
		dsn,
		otelsql.WithAttributes(attrs...),
		otelsql.WithTextMapPropagator(params.Propagator),
		otelsql.WithMeterProvider(params.MeterProvider),
		otelsql.WithTracerProvider(params.TracerProvider),
	)
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

type ProvideMigratorParams struct {
	fx.In
	Config  *config.Config
	Logger  *slog.Logger
	DB      *sql.DB
	Options []MigratorOption `group:"db-migrator-options"`
}

func ProvideMigrator(params ProvideMigratorParams) *Migrator {
	return NewMigrator(params.Logger, params.DB, params.Options...)
}

func RunMigrations(command string, args ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, migrator *Migrator, config *config.Config) error {
			return migrator.Run(ctx, config.GetString("db.driver"), command, args...)
		},
	)
}

func RunMigrationsAndShutdown(command string, args ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, migrator *Migrator, config *config.Config, shutdown fx.Shutdowner) error {
			defer shutdown.Shutdown()

			return migrator.Run(ctx, config.GetString("db.driver"), command, args...)
		},
	)
}

type ProvideSeederParams struct {
	fx.In
	Config *config.Config
	Logger *slog.Logger
	DB     *sql.DB
	Seeds  []Seed `group:"db-seeder-seeds"`
}

func ProvideSeeder(params ProvideSeederParams) *Seeder {
	return NewSeeder(params.Logger, params.DB, params.Seeds...)
}

func RunSeeds(names ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, seeder *Seeder) error {
			return seeder.Run(ctx, names...)
		},
	)
}

func RunSeedsAndShutdown(names ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, seeder *Seeder, shutdown fx.Shutdowner) error {
			defer shutdown.Shutdown()

			return seeder.Run(ctx, names...)
		},
	)
}
