package test

import (
	"context"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/tingold/bunpostgis"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"log"
	"os"
	"testing"
)

var database *bun.DB

func TestMain(m *testing.M) {

	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	db := os.Getenv("PGDATABASE")
	user := os.Getenv("PGUSER")
	pw := os.Getenv("PGPASSWORD")

	connstring := "database=" + db + " host=" + host + " user=" + user + " password=" + pw + " port=" + port
	config, err := pgx.ParseConfig(connstring)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	sqldb := stdlib.OpenDB(*config)
	database = bun.NewDB(sqldb, pgdialect.New())
	defer database.Close()

	err = database.ResetModel(context.Background(), &TestStruct{})
	if err != nil {
		log.Fatalf("unable to reset model :%s", err.Error())
	}
	m.Run()

}

type TestStruct struct {
	bun.Model `bun:"bunpostgis_test_table"`
	GeoField  bunpostgis.PostgisGeometry `bun:geom`
}
