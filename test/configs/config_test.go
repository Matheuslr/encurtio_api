package configs_test

import (
	"os"
	"testing"

	"github.com/matheuslr/encurtio/configs"
	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("CASSANDRA_HOST")
	os.Unsetenv("CASSANDRA_KEYSPACE")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("APP_URL")

	cfg := configs.Load()

	assert.NotNil(t, cfg)
	assert.Equal(t, []string{"127.0.0.1"}, cfg.Cassandra.Hosts)
	assert.Equal(t, "encurtio", cfg.Cassandra.Keyspace)
	assert.Equal(t, "8080", cfg.API.Port)
	assert.Equal(t, "http://localhost:8080", cfg.API.URL)
}

func TestLoad_FromEnvVars(t *testing.T) {
	os.Setenv("CASSANDRA_HOST", "192.168.1.100")
	os.Setenv("CASSANDRA_KEYSPACE", "test_keyspace")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("APP_URL", "https://short.io")
	defer func() {
		os.Unsetenv("CASSANDRA_HOST")
		os.Unsetenv("CASSANDRA_KEYSPACE")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_URL")
	}()

	cfg := configs.Load()

	assert.Equal(t, []string{"192.168.1.100"}, cfg.Cassandra.Hosts)
	assert.Equal(t, "test_keyspace", cfg.Cassandra.Keyspace)
	assert.Equal(t, "9090", cfg.API.Port)
	assert.Equal(t, "https://short.io", cfg.API.URL)
}
