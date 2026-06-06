package geo

import (
	"math"
)

const (
	MinLat = -85.05112878
	MaxLat = 85.05112878
	MinLon = -180.0
	MaxLon = 180.0

	latRange = MaxLat - MinLat
	lonRange = MaxLon - MinLon
)

func compactInt64ToInt32(v uint64) uint32 {
	result := v & 0x5555555555555555
	result = (result | (result >> 1)) & 0x3333333333333333
	result = (result | (result >> 2)) & 0x0F0F0F0F0F0F0F0F
	result = (result | (result >> 4)) & 0x00FF00FF00FF00FF
	result = (result | (result >> 8)) & 0x0000FFFF0000FFFF
	result = (result | (result >> 16)) & 0x00000000FFFFFFFF
	return uint32(result)
}

func convertGridNumbersToCoordinates(gridLatitudeNumber, gridLongitudeNumber uint32) Coord {
	// Calculate the grid boundaries
	gridLatitudeMin := MinLat + latRange*(float64(gridLatitudeNumber)/math.Pow(2, 26))
	gridLatitudeMax := MinLat + latRange*(float64(gridLatitudeNumber+1)/math.Pow(2, 26))
	gridLongitudeMin := MinLon + lonRange*(float64(gridLongitudeNumber)/math.Pow(2, 26))
	gridLongitudeMax := MinLon + lonRange*(float64(gridLongitudeNumber+1)/math.Pow(2, 26))

	// Calculate the center point of the grid cell
	latitude := (gridLatitudeMin + gridLatitudeMax) / 2
	longitude := (gridLongitudeMin + gridLongitudeMax) / 2

	return Coord{Lat: latitude, Lon: longitude}
}

func decode(geoCode uint64) Coord {
	// Align bits of both latitude and longitude to take even-numbered position
	y := geoCode >> 1
	x := geoCode

	// Compact bits back to 32-bit ints
	gridLatitudeNumber := compactInt64ToInt32(x)
	gridLongitudeNumber := compactInt64ToInt32(y)

	return convertGridNumbersToCoordinates(gridLatitudeNumber, gridLongitudeNumber)
}
