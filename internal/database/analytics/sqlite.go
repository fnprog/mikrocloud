package analytics_db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"
)

type SQLiteAnalyticsDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(connectionString string) (AnalyticsDatabase, error) {
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Analytics database connection established (SQLite)", "connection", connectionString)

	instance := &SQLiteAnalyticsDatabase{db: db}

	if err := instance.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return instance, nil
}

func (s *SQLiteAnalyticsDatabase) Close() error {
	if err := s.db.Close(); err != nil {
		slog.Error("Error closing analytics database", "error", err)
		return err
	}
	slog.Info("Analytics database connection closed")
	return nil
}

func (s *SQLiteAnalyticsDatabase) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *SQLiteAnalyticsDatabase) DB() interface{} {
	return s.db
}

func (s *SQLiteAnalyticsDatabase) initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS metrics (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			service_id TEXT,
			name TEXT NOT NULL,
			value REAL NOT NULL,
			unit TEXT NOT NULL,
			tags TEXT NOT NULL DEFAULT '{}',
			timestamp BIGINT NOT NULL,
			created_at BIGINT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			source TEXT NOT NULL,
			data TEXT,
			tags TEXT,
			timestamp DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id TEXT PRIMARY KEY,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			source TEXT NOT NULL,
			fields TEXT,
			tags TEXT,
			timestamp DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_project_id ON metrics(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics(name)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_name_timestamp ON metrics(name, timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type_timestamp ON events(type, timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_level_timestamp ON logs(level, timestamp)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute schema query: %w", err)
		}
	}

	return nil
}

func (s *SQLiteAnalyticsDatabase) StoreMetric(ctx context.Context, metric Metric) error {
	tagsJSON, err := json.Marshal(metric.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `INSERT INTO metrics (id, name, value, tags, timestamp) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, metric.ID, metric.Name, metric.Value, string(tagsJSON), metric.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to store metric: %w", err)
	}

	return nil
}

func (s *SQLiteAnalyticsDatabase) StoreEvent(ctx context.Context, event Event) error {
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	tagsJSON, err := json.Marshal(event.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `INSERT INTO events (id, type, source, data, tags, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, event.ID, event.Type, event.Source, string(dataJSON), string(tagsJSON), event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	return nil
}

func (s *SQLiteAnalyticsDatabase) StoreLogs(ctx context.Context, logs []LogEntry) error {
	if len(logs) == 0 {
		return nil
	}

	query := `INSERT INTO logs (id, level, message, source, fields, tags, timestamp) VALUES `
	values := make([]string, len(logs))
	args := make([]any, 0, len(logs)*7)

	for i, log := range logs {
		fieldsJSON, err := json.Marshal(log.Fields)
		if err != nil {
			return fmt.Errorf("failed to marshal fields for log %s: %w", log.ID, err)
		}

		tagsJSON, err := json.Marshal(log.Tags)
		if err != nil {
			return fmt.Errorf("failed to marshal tags for log %s: %w", log.ID, err)
		}

		values[i] = "(?, ?, ?, ?, ?, ?, ?)"
		args = append(args, log.ID, log.Level, log.Message, log.Source, string(fieldsJSON), string(tagsJSON), log.Timestamp)
	}

	query += strings.Join(values, ", ")
	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to store logs: %w", err)
	}

	return nil
}

func (s *SQLiteAnalyticsDatabase) GetMetrics(ctx context.Context, filter MetricFilter) ([]Metric, error) {
	query := `SELECT id, name, value, tags, timestamp FROM metrics WHERE timestamp >= ? AND timestamp <= ?`
	args := []any{filter.TimeRange.Start, filter.TimeRange.End}

	if len(filter.Names) > 0 {
		placeholders := make([]string, len(filter.Names))
		for i, name := range filter.Names {
			placeholders[i] = "?"
			args = append(args, name)
		}
		query += ` AND name IN (` + strings.Join(placeholders, ",") + `)`
	}

	query += ` ORDER BY timestamp DESC`
	if filter.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, filter.Limit)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var metric Metric
		var tagsJSON string

		if err := rows.Scan(&metric.ID, &metric.Name, &metric.Value, &tagsJSON, &metric.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %w", err)
		}

		if err := json.Unmarshal([]byte(tagsJSON), &metric.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (s *SQLiteAnalyticsDatabase) GetEvents(ctx context.Context, filter EventFilter) ([]Event, error) {
	return []Event{}, nil
}

func (s *SQLiteAnalyticsDatabase) GetLogs(ctx context.Context, filter LogFilter) ([]LogEntry, error) {
	return []LogEntry{}, nil
}

func (s *SQLiteAnalyticsDatabase) GetMetricAggregation(ctx context.Context, aggregation MetricAggregation) ([]AggregationResult, error) {
	return []AggregationResult{}, nil
}

func (s *SQLiteAnalyticsDatabase) GetEventStats(ctx context.Context, timeRange TimeRange) (EventStats, error) {
	return EventStats{
		TimeRange:      timeRange,
		EventsByType:   make(map[string]int64),
		EventsBySource: make(map[string]int64),
	}, nil
}
