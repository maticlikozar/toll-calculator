-- Create events table
CREATE TABLE IF NOT EXISTS events (
    created_at TIMESTAMPTZ NOT NULL,
    license_plate TEXT NOT NULL,
    event_start TIMESTAMPTZ NOT NULL,
    event_stop TIMESTAMPTZ NOT NULL,
    car_type TEXT NOT NULL,
    billed BOOLEAN NOT NULL DEFAULT FALSE
);

SELECT create_hypertable('events', 'created_at', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS api_key (
    id          BYTEA NOT NULL,
    key_hash    VARCHAR(256) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT unique_key UNIQUE (key_hash)
);
