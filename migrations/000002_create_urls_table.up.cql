CREATE TABLE IF NOT EXISTS encurtio.urls (
    short_code text PRIMARY KEY,
    original_url text,
    created_at timestamp
);