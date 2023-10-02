package scylla

import (
	"fmt"
	"time"

	"github.com/BigBoulard/scylla-mutants/internal/logger"
	"github.com/gocql/gocql"
)

func CreateCluster(consistency gocql.Consistency, keyspace string, hosts ...string) *gocql.ClusterConfig {
	retryPolicy := &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        10 * time.Second,
		NumRetries: 5,
	}
	config := gocql.NewCluster(hosts...)
	config.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	config.Compressor = &gocql.SnappyCompressor{}
	config.RetryPolicy = retryPolicy
	config.Timeout = 5 * time.Second

	// config.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: conf.DB.User,
	// 	Password: conf.DB.Pass,
	// }

	config.Keyspace = keyspace
	config.Consistency = consistency

	return config
}

// see https://university.scylladb.com/courses/using-scylla-drivers/lessons/golang-and-scylla-part-1/
// see https://github.com/scylladb/scylla-code-samples/tree/master/mms/go
// see https://github.com/scylladb/gocql
func Start(log logger.Logger) {
	// clusterConf := gocql.NewCluster(conf.DB.Host)

	// config := CreateCluster(gocql.Quorum, "catalog", "scylla-node1", "scylla-node2", "scylla-node3")
	config := CreateCluster(gocql.Quorum, "catalog", "172.29.0.5", "172.29.0.4", "172.29.0.3")
	session, err := gocql.NewSession(*config)
	if err != nil {
		log.Fatal(err, "db", "Start", "unable to connect to scylla")
	}
	defer session.Close()
}

func SelectQuery(session *gocql.Session, log logger.Logger) {
	log.Info("db", "SelectQuery", "Displaying Results")
	q := session.Query("SELECT first_name,last_name,address,picture_location FROM mutant_data")
	var firstName, lastName, address, pictureLocation string
	it := q.Iter()
	defer func() {
		if err := it.Close(); err != nil {
			log.Warn("db", "SelectQuery", fmt.Sprintf("select catalog.mutant error: %s", err))
		}
	}()
	for it.Scan(&firstName, &lastName, &address, &pictureLocation) {
		log.Info("db", "SelectQuery", "\t"+firstName+" "+lastName+", "+address+", "+pictureLocation)
	}
}

func InsertQuery(session *gocql.Session, log logger.Logger) {
	log.Info("db", "SelectQuery", "Inserting Mike")
	if err := session.Query("INSERT INTO mutant_data (first_name,last_name,address,picture_location) VALUES ('Mike','Tyson','1515 Main St', 'http://www.facebook.com/mtyson')").Exec(); err != nil {
		log.Error(err, "db", "SelectQuery", "insert catalog.mutant_data")
	}
}

func DeleteQuery(session *gocql.Session, log logger.Logger) {
	log.Info("db", "SelectQuery", "Deleting Mike")
	if err := session.Query("DELETE FROM mutant_data WHERE first_name = 'Mike' and last_name = 'Tyson'").Exec(); err != nil {
		log.Error(err, "db", "SelectQuery", "delete catalog.mutant_data")
	}
}
