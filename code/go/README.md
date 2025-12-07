# Local setup

Start docker environment and make sure _toll-events_ and _toll-api_ containers are running.

```bash
$ cd backend && docker-compose up toll-api toll-events -d
```

making request cal to /health endpoint should return 200 OK response.

```curl
$ curl --location --request GET 'http://localhost:8080/api/v1/health'
```

## Setup local API key in timescaleDB

Connect to timescaleDB docker container and create new _api_key_ table

```sql
CREATE TABLE IF NOT EXISTS api_key (
    id          BYTEA NOT NULL,
    key_hash    VARCHAR(256) NOT NULL,
    system_key  BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at  TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT unique_key UNIQUE (key_hash)
);
```

Insert new api key in _api_key_ table.

```sql
INSERT INTO api_key (
    id,
    key_hash,
    expires_at
) VALUES (
    decode(replace(gen_random_uuid()::text, '-', ''), 'hex'),
    'pmWkWSBCL51Bfkhn79xPuKBKHz__H6B-mY6G9_eieuM',
    NOW() + INTERVAL '30 days'
);
```

Use created api_key as X-API-Key header when requesting api.

```curl
curl --location 'https://toll.test/api/v1/toll-events' \
--header 'X-API-Key: 123' \
--header 'Content-Type: application/json' \
--data '{
    ...
}'
```
