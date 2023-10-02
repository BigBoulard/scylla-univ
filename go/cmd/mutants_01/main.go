package main

import (
	"github.com/BigBoulard/scylla-mutants/internal/logger"
	"github.com/BigBoulard/scylla-mutants/internal/scylla"
	"github.com/gocql/gocql"
)

func main() {
	logger := logger.NewLogger()
	clusterConf := scylla.CreateCluster(gocql.Quorum, "catalog", "scylla-node1", "scylla-node2", "scylla-node3")
	session, err := gocql.NewSession(*clusterConf)

	if err != nil {
		logger.Fatal(err, "main", "main", "unable to connect to scylla")
	}
	defer session.Close()

	scylla.SelectQuery(session, logger)
	scylla.InsertQuery(session, logger)
	scylla.SelectQuery(session, logger)
	scylla.DeleteQuery(session, logger)
	scylla.SelectQuery(session, logger)
}
