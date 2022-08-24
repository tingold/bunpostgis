# PostGIS support for Bun ORM

[![Build, Lint, Test](https://github.com/tingold/bunpostgis/actions/workflows/go.yml/badge.svg)](https://github.com/tingold/bunpostgis/actions/workflows/go.yml)  [![Go Report Card](https://goreportcard.com/badge/github.com/tingold/bunpostgis/pkg)](https://goreportcard.com/report/github.com/tingold/bunpostgis/pkg)

This module provides a simple wrapper that allows [PostGIS](https://postgis.net/) geometry types to be used in [Bun](https://bun.uptrace.dev/).
It builds on the awesome `ewkb` support already in [Orb](https://github.com/paulmach/orb). 

**Only pgx supported**

Supports Point/MultiPoint, Linestring/MultiLinestring, Polygon/MultiPolygon.

### Usage

Use it in a struct and tag the type appropriately if you want Bun to be able to create the table:

```
type SampleStruct struct {
	bun.BaseModel `bun:"bunpostgis_test_table"`
	GeoField      bunpostgis.PostgisGeometry `bun:"type:Geometry"`
	Name          string
}
```

See [`test/postgis_test.go`](test/postgis_test.go) for a full example 

## Testing

Tests have been refactored into their own package to avoid unessesary imports when using the library.

This test expects [Postgres environment variables](https://www.postgresql.org/docs/current/libpq-envars.html) to be present to configure the database connection -- intended mostly for CI. 
