package test

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/paulmach/orb"
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

	bunDb := bun.NewDB(stdlib.OpenDB(*config), pgdialect.New())
	database = bunDb
	_, err = bunDb.NewCreateTable().Model((*SampleStruct)(nil)).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalf("unable to reset model :%s", err.Error())
	}
	m.Run()
	//cleanup
	_, err = bunDb.NewDropTable().Model((*SampleStruct)(nil)).IfExists().Exec(context.Background())
	if err != nil {
		log.Fatalf("unable to cleanup test table: %s", err.Error())
	}
}

func TestAllGeometries(t *testing.T) {

	geoms := []bunpostgis.PostgisGeometry{
		bunpostgis.PostgisGeometry{Geometry: orb.Point{-76.35, 39.53}, SRID: 4326},
		bunpostgis.PostgisGeometry{Geometry: orb.Polygon{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}}, SRID: 4326},
		bunpostgis.PostgisGeometry{Geometry: orb.LineString{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}, SRID: 4326},
		bunpostgis.PostgisGeometry{Geometry: orb.MultiPoint{{-76.35, 39.53}, {22, 12}}, SRID: 4326},
		bunpostgis.PostgisGeometry{Geometry: orb.MultiPolygon{
			{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}},
			{{{1, 1}, {1, 3}, {3, 3}, {3, 1}, {1, 1}}},
		}, SRID: 4326},
		bunpostgis.PostgisGeometry{Geometry: orb.MultiLineString{
			{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}},
			{{1, 1}, {1, 3}, {3, 3}, {3, 1}, {1, 1}},
		}, SRID: 4326},
	}

	ctx := context.Background()
	for k, v := range geoms {

		s := SampleStruct{GeoField: v, Name: fmt.Sprintf("geometry %d", k)}
		_, err := database.NewInsert().Model(&s).Exec(ctx)
		if err != nil {
			t.Errorf("error inserting struct: %s", err.Error())
			t.Fail()
		}
		var s2 SampleStruct
		err = database.NewSelect().Model(&s2).Where("name=?", s.Name).Scan(ctx, &s2)
		if err != nil {
			t.Errorf("error loading struct: %s", err.Error())
			t.Fail()
		}
		if s2.Name != s.Name {
			t.Fail()
		}
		if s2.GeoField.SRID != s.GeoField.SRID {
			t.Fail()
		}
		if s2.GeoField.Geometry.GeoJSONType() != s.GeoField.Geometry.GeoJSONType() {
			t.Fail()
		}
		if s2.GeoField.Geometry.Bound() != s.GeoField.Geometry.Bound() {
			t.Fail()
		}
	}

}

type SampleStruct struct {
	bun.BaseModel `bun:"bunpostgis_test_table"`
	GeoField      bunpostgis.PostgisGeometry `bun:"type:Geometry"`
	Name          string
}
