CREATE KEYSPACE IF NOT EXISTS encurtio
WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};