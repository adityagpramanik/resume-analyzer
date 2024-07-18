package commonservices

import (
	"log"
	"os"
	"strings"

	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func getSession() *gocql.Session {
	hosts := os.Getenv("CASSANDRA_HOSTS")
	if hosts == "" {
		log.Fatalf("CASSANDRA_HOSTS environment variable is required")
	}

	hostList := strings.Split(hosts, ",")

	cluster := gocql.NewCluster(hostList...)
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Two

	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	return session
}

func GetSession() *gocql.Session {
	return getSession()
}

func CloseSession() {
	if session != nil {
		session.Close()
	}
}
