package main

import (
	"fmt"

	"github.com/BigBoulard/scylla-mutants/internal/logger"
	"github.com/BigBoulard/scylla-mutants/internal/scylla"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
)

var stmts = createStatements()

func main() {
	logger := logger.NewLogger()

	cluster := scylla.CreateCluster(gocql.Quorum, "catalog", "scylla-node1", "scylla-node2", "scylla-node3")
	session, err := gocql.NewSession(*cluster)
	if err != nil {
		logger.Fatal(err, "mutants_03", "", "unable to connect to scylla")
	}
	defer session.Close()

	selectQuery(session, logger)
	insertQuery(session, "Mike", "Tyson", "12345 Foo Lane", "http://www.facebook.com/mtyson", logger)
	insertQuery(session, "Alex", "Jones", "56789 Hickory St", "http://www.facebook.com/ajones", logger)
	selectQuery(session, logger)
	deleteQuery(session, "Mike", "Tyson", logger)
	selectQuery(session, logger)
	deleteQuery(session, "Alex", "Jones", logger)
	selectQuery(session, logger)
}

func deleteQuery(session *gocql.Session, firstName string, lastName string, logger logger.Logger) {
	logger.Info("mutants_03", "insertQuery", "Deleting "+firstName+"......")
	r := Record{
		FirstName: firstName,
		LastName:  lastName,
	}
	err := gocqlx.Query(session.Query(stmts.del.stmt), stmts.del.names).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error(err, "mutants_03", "deleteQuery", "delete catalog.mutant_data")
	}
}

func insertQuery(session *gocql.Session, firstName, lastName, address, pictureLocation string, logger logger.Logger) {
	logger.Info("mutants_03", "insertQuery", "Inserting "+firstName+"......")
	r := Record{
		FirstName:       firstName,
		LastName:        lastName,
		Address:         address,
		PictureLocation: pictureLocation,
	}
	err := gocqlx.Query(session.Query(stmts.ins.stmt), stmts.ins.names).BindStruct(r).ExecRelease()
	if err != nil {
		logger.Error(err, "mutants_03", "insertQuery", "insert catalog.mutant_data")
	}
}

func selectQuery(session *gocql.Session, logger logger.Logger) {
	logger.Info("mutants_03", "insertQuery", "Displaying Results:")
	var rs []Record
	err := gocqlx.Query(session.Query(stmts.sel.stmt), stmts.sel.names).SelectRelease(&rs)
	if err != nil {
		logger.Warn("mutants_03", "selectQuery", fmt.Sprintf("select catalog.mutant %v", err))
		return
	}
	for _, r := range rs {
		logger.Info("mutants_03", "insertQuery", "\t"+r.FirstName+" "+r.LastName+", "+r.Address+", "+r.PictureLocation)
	}
}

func createStatements() *statements {
	m := table.Metadata{
		Name:    "mutant_data",
		Columns: []string{"first_name", "last_name", "address", "picture_location"},
		PartKey: []string{"first_name", "last_name"},
	}
	tbl := table.New(m)

	deleteStmt, deleteNames := tbl.Delete()
	insertStmt, insertNames := tbl.Insert()
	// Normally a select statement such as this would use `tbl.Select()` to select by
	// primary key but now we just want to display all the records...
	selectStmt, selectNames := qb.Select(m.Name).Columns(m.Columns...).ToCql()
	return &statements{
		del: query{
			stmt:  deleteStmt,
			names: deleteNames,
		},
		ins: query{
			stmt:  insertStmt,
			names: insertNames,
		},
		sel: query{
			stmt:  selectStmt,
			names: selectNames,
		},
	}
}

type query struct {
	stmt  string
	names []string
}

type statements struct {
	del query
	ins query
	sel query
}

type Record struct {
	FirstName       string `db:"first_name"`
	LastName        string `db:"last_name"`
	Address         string `db:"address"`
	PictureLocation string `db:"picture_location"`
}
