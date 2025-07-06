CREATE TABLE IF NOT EXISTS
    sensors.birdnet (
        time TIMESTAMPTZ NOT NULL,
        source_node TEXT NOT NULL,
        source TEXT NOT NULL,
        begin_time TIMESTAMPTZ NOT NULL,
        end_time TIMESTAMPTZ NOT NULL,
        species_code TEXT NOT NULL,
        scientific_name TEXT NOT NULL,
        common_name TEXT NOT NULL,
        confidence FLOAT8 NOT NULL,
        latitude FLOAT8 NOT NULL,
        longitude FLOAT8 NOT NULL,
        threshold FLOAT8 NOT NULL,
        sensitivity FLOAT8 NOT NULL
    );

SELECT
    create_hypertable ('sensors.birdnet', 'time', if_not_exists => TRUE);