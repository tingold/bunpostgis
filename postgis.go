package bunpostgis

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
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

	switch value := value.(type) {
	case string:
		data, err = hex.DecodeString(value)
		if err != nil {
			return err
		}

		g.Geometry, g.SRID, err = ewkb.Unmarshal(data)
		return err
	default:
		return fmt.Errorf("expected string but got %T", value)
	}
}

func (g *PostgisGeometry) Value() (driver.Value, error) {
	if g.Geometry == nil {
		return nil, nil
	}

	d := ewkb.MustMarshalToHex(g.Geometry, g.SRID)
	return d, nil

}
