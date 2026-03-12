package database

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

type CassandraConfig struct {
	Hosts    []string
	Keyspace string
}

func NewCassandraSession(cfg CassandraConfig) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)

	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Cassandra")

	return session, nil

}
