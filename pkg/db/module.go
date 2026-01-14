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
