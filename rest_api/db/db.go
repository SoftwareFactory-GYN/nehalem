package db

import (
	"github.com/gocql/gocql"
)

func GetSession() *gocql.Session {

	var err error

	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "nehalem"
	Session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	return Session
}
