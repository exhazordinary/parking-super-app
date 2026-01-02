package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	dbTracerName = "database"
)

// DBSpanConfig holds configuration for database span creation
type DBSpanConfig struct {
	DBSystem    string // e.g., "postgresql"
	DBName      string
	DBUser      string
	ServerAddr  string
	ServerPort  int
}

// TraceDBQuery creates a span for a database query
// Returns the context with the span and the span itself for manual ending
func TraceDBQuery(ctx context.Context, operation, query string, cfg *DBSpanConfig) (context.Context, trace.Span) {
	tracer := otel.Tracer(dbTracerName)

	attrs := []attribute.KeyValue{
		attribute.String("db.operation", operation),
	}

	if cfg != nil {
		if cfg.DBSystem != "" {
			attrs = append(attrs, attribute.String("db.system", cfg.DBSystem))
		}
		if cfg.DBName != "" {
			attrs = append(attrs, attribute.String("db.name", cfg.DBName))
		}
		if cfg.DBUser != "" {
			attrs = append(attrs, attribute.String("db.user", cfg.DBUser))
		}
		if cfg.ServerAddr != "" {
			attrs = append(attrs, attribute.String("server.address", cfg.ServerAddr))
		}
		if cfg.ServerPort > 0 {
			attrs = append(attrs, attribute.Int("server.port", cfg.ServerPort))
		}
	}

	// Include query statement (consider truncating for very long queries)
	if query != "" {
		if len(query) > 1000 {
			query = query[:1000] + "..."
		}
		attrs = append(attrs, attribute.String("db.statement", query))
	}

	ctx, span := tracer.Start(ctx, "db."+operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)

	return ctx, span
}

// TraceDBQueryFunc is a helper for tracing database operations with automatic span management
func TraceDBQueryFunc[T any](ctx context.Context, operation, query string, cfg *DBSpanConfig, fn func(context.Context) (T, error)) (T, error) {
	ctx, span := TraceDBQuery(ctx, operation, query, cfg)
	defer span.End()

	result, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return result, err
}

// TraceDBExec is a helper for tracing database exec operations (INSERT, UPDATE, DELETE)
func TraceDBExec(ctx context.Context, operation, query string, cfg *DBSpanConfig, fn func(context.Context) error) error {
	ctx, span := TraceDBQuery(ctx, operation, query, cfg)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}
