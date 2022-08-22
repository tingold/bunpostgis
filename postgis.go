package bunpostgis

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"reflect"
)

type PostgisGeometry struct {
	orb.Geometry
	SRID int
}

func (g *PostgisGeometry) Scan(value interface{}) error {

	if value == nil {
		return nil
	}
	var data []byte
	var err error
	if reflect.TypeOf(value).Kind() == reflect.String {
		data, err = hex.DecodeString(value.(string))
		if err != nil {
			return err
		}
	} else {
		//need to throw some error here
		return fmt.Errorf("expected string but got %s", reflect.TypeOf(value).Kind().String())
	}

	g.Geometry, g.SRID, err = ewkb.Unmarshal(data)
	return err
}

func (g *PostgisGeometry) Value() (driver.Value, error) {

	if g.Geometry == nil {
		return nil, nil
	}
	d := ewkb.MustMarshalToHex(g.Geometry, g.SRID)
	return d, nil

}
