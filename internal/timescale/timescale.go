package timescale

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelpeterswa/birdnet-mqtt-timescaledb/internal/birdnet"
)

//go:embed queries/insert_bird_detection.sql
var insertBirdDetectionQuery string

type TimescaleClient struct {
	Pool *pgxpool.Pool
}

func NewTimescaleClient(ctx context.Context, connString string) (*TimescaleClient, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &TimescaleClient{Pool: pool}, nil
}

func (c *TimescaleClient) Close() {
	c.Pool.Close()
}

func (c *TimescaleClient) StoreBirdDetectionEvent(ctx context.Context, event *birdnet.BirdDetectionEvent) error {
	_, err := c.Pool.Exec(ctx, insertBirdDetectionQuery,
		event.Time, event.SourceNode, event.Source,
		event.BeginTime, event.EndTime, event.SpeciesCode,
		event.ScientificName, event.CommonName, event.Confidence,
		event.Latitude, event.Longitude, event.Threshold, event.Sensitivity,
	)

	return err
}
