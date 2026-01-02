package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Tracing returns HTTP middleware that adds OpenTelemetry tracing
func Tracing(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming request headers
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Create span name from method and path
			spanName := r.Method + " " + r.URL.Path

			// Start span
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethod(r.Method),
					semconv.HTTPTarget(r.URL.Path),
					semconv.HTTPScheme(getScheme(r)),
					attribute.String("http.host", r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
				),
			)
			defer span.End()

			// Wrap response writer to capture status code
			rw := newResponseWriter(w)

			// Call next handler with traced context
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Record response attributes
			span.SetAttributes(
				semconv.HTTPStatusCode(rw.statusCode),
			)

			// Set span status based on HTTP status code
			if rw.statusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(rw.statusCode))
			} else {
				span.SetStatus(codes.Ok, "")
			}
		})
	}
}

// TracingWithConfig returns HTTP middleware with custom configuration
type TracingConfig struct {
	ServiceName    string
	SkipPaths      []string // Paths to skip tracing (e.g., /health)
	RecordHeaders  bool     // Whether to record request headers
	PropagateTrace bool     // Whether to propagate trace context
}

func TracingWithConfig(cfg TracingConfig) func(http.Handler) http.Handler {
	tracer := otel.Tracer(cfg.ServiceName)
	propagator := otel.GetTextMapPropagator()

	skipPaths := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip tracing for certain paths
			if skipPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			// Extract trace context
			ctx := r.Context()
			if cfg.PropagateTrace {
				ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
			}

			spanName := r.Method + " " + r.URL.Path

			attrs := []attribute.KeyValue{
				semconv.HTTPMethod(r.Method),
				semconv.HTTPTarget(r.URL.Path),
				semconv.HTTPScheme(getScheme(r)),
				attribute.String("http.host", r.Host),
			}

			if cfg.RecordHeaders {
				for key, values := range r.Header {
					if len(values) > 0 {
						attrs = append(attrs, attribute.String("http.request.header."+key, values[0]))
					}
				}
			}

			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(attrs...),
			)
			defer span.End()

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r.WithContext(ctx))

			span.SetAttributes(semconv.HTTPStatusCode(rw.statusCode))
			if rw.statusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(rw.statusCode))
			}
		})
	}
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}
