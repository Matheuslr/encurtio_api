CREATE KEYSPACE IF NOT EXISTS encurtio
WITH replication = {
 'class': 'SimpleStrategy',
 'replication_factor': 1
};

USE encurtio;

CREATE TABLE IF NOT EXISTS urls (
    short_code text PRIMARY KEY,
    id bigint,
    original_url text,
    created_at timestamp
);