CREATE DATABASE IF NOT EXISTS logs;

CREATE TABLE IF NOT EXISTS logs.app_logs
(
    timestamp DateTime64(3, 'UTC'),

    level LowCardinality(String),

    message String,

    request_id String DEFAULT '',

    service LowCardinality(String) DEFAULT 'tmaster',

    container_name String DEFAULT '',

    attributes JSON
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, level);