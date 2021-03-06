package gosm
// Forked from github.com/apeyroux/gosm
// In order to add necessary minor fixes

import (
	"math"
)

type Tile struct {
	Z    int64
	X    int64
	Y    int64
	Lat  float64
	Long float64
}

type BBox struct {
	TopLeftTile     Tile
	BottomRightTile Tile
}

type Conversion interface {
	Deg2num(t *Tile) (x int64, y int64)
	Num2deg(t *Tile) (lat float64, long float64)
}

func (t *Tile) Deg2num() (x int64, y int64) {
	x = int64(math.Floor((t.Long + 180.0) / 360.0 * (math.Exp2(float64(t.Z)))))
	y = int64(math.Floor((1.0 - math.Log(math.Tan(t.Lat*math.Pi/180.0)+1.0/math.Cos(t.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(t.Z)))))
	return
}

func (t *Tile) Num2deg() (lat float64, long float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return lat, long
}

func NewTileWithLatLong(lat float64, long float64, z int64) (t *Tile) {
	t = new(Tile)
	t.Lat = lat
	t.Long = long
	t.Z = z
	t.X, t.Y = t.Deg2num()
	return
}

func NewTileWithXY(x int64, y int64, z int64) (t *Tile) {
	t = new(Tile)
	t.Z = z
	t.X = x
	t.Y = y
	t.Lat, t.Long = t.Num2deg()
	return
}

func BBoxTiles(topTile Tile, bottomTile Tile) ([]*Tile, error) {
	tiles := []*Tile{}
	for _, z := range []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19} {
		tMax := NewTileWithLatLong(topTile.Lat, topTile.Long, z)
		tMin := NewTileWithLatLong(bottomTile.Lat, bottomTile.Long, z)
		//nbtiles := math.Abs((float64(tMax.X))-float64(tMin.X)) + math.Abs(float64(tMax.Y)-float64(tMin.Y))
		for x := tMin.X; x <= tMax.X; x++ {
			for y := tMax.Y; y <= tMin.Y; y++ {
				tiles = append(tiles, NewTileWithXY(x, y, z))
			}
		}
	}
	return tiles, nil
}
