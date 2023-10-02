package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BigBoulard/scylla-mutants/internal/scylla"

	"github.com/BigBoulard/scylla-mutants/internal/logger"
	"github.com/gocql/gocql"
	"github.com/qeesung/image2ascii/convert"
)

var (
	alter bool
)

func init() {
	flag.BoolVar(&alter, "alter", false, "issue alter statements")
}

func main() {
	flag.Parse()

	logger := logger.NewLogger()

	cluster := scylla.CreateCluster(gocql.Quorum, "catalog", "scylla-node1", "scylla-node2", "scylla-node3")
	session, err := gocql.NewSession(*cluster)
	if err != nil {
		logger.Fatal(err, "mutants_02", "main", "unable to connect to scylla")
	}
	defer session.Close()

	if alter {
		alterSchema(session, logger)
	}

	mutants := []string{"Jim Jefferies", "Bob Loblaw", "Bob Zemuda"}
	for _, mutant := range mutants {
		names := strings.SplitN(mutant, " ", 2)
		firstName := names[0]
		lastName := names[1]
		logger.Info("mutants_02", "main", "Processing file for "+firstName+"_"+lastName)
		insertFile(session, firstName, lastName, logger)
		readFile(session, firstName, lastName, logger)
	}
}

func readFile(session *gocql.Session, firstName string, lastName string, logger logger.Logger) {
	var (
		m map[string][]byte
		b []byte
	)
	if err := session.Query("SELECT m,b from mutant_data WHERE first_name=? AND last_name=?", firstName, lastName).Scan(&m, &b); err != nil {
		logger.Fatal(err, "mutants_02", "readFile", fmt.Sprintf("unable to read image of %s %s", firstName, lastName))
	}
	logger.Info("mutants_02", "readFile", fmt.Sprintf("file metadata: %+v", m))

	if err := os.WriteFile("/tmp/"+firstName+"_"+lastName+".jpg", b, 0644); err != nil {
		logger.Fatal(err, "mutants_02", "readFile", fmt.Sprintf("unable to write image of %s %s in %s", firstName, lastName, "/tmp/"+firstName+"_"+lastName+".jpg"))
	}
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 100
	convertOptions.FixedHeight = 40

	converter := convert.NewImageConverter()
	fmt.Print(converter.ImageFile2ASCIIString("/tmp/"+firstName+"_"+lastName+".jpg", &convertOptions))
}

func alterSchema(session *gocql.Session, logger logger.Logger) {
	if err := session.Query("ALTER table catalog.mutant_data ADD m map<text, blob>").Exec(); err != nil {
		logger.Fatal(err, "mutants_02", "alterSchema", "altering table failed, forgot the '-alter=false' program argument?")
	}
	if err := session.Query("ALTER table catalog.mutant_data ADD b blob").Exec(); err != nil {
		logger.Fatal(err, "mutants_02", "alterSchema", "altering table failed, forgot the '-alter=false' program argument?")
	}
}

func insertFile(session *gocql.Session, firstName, lastName string, logger logger.Logger) {
	fName := "/usr/share/icons/mutants/" + firstName + "_" + lastName + ".jpg"
	b, err := os.ReadFile(fName)
	if err != nil {
		logger.Fatal(err, "mutants_02", "insertFile", fmt.Sprintf("unable to read file %s", fName))
	}

	m := readMetaFromFile(firstName+"_"+lastName, b)

	if err := session.Query("INSERT INTO mutant_data (first_name,last_name,b,m) VALUES (?,?,?,?)", firstName, lastName, b, m).Exec(); err != nil {
		logger.Error(err, "mutants_02", "insertFile", "insert catalog.mutant_data error")
	}
}

func readMetaFromFile(name string, bytes []byte) map[string][]byte {
	return map[string][]byte{
		"name": []byte(name),
	}
}
