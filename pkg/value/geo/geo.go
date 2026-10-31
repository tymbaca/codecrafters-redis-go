package geo

type Coord struct {
	Lon, Lat float64
}

func FromScore(score float64) Coord {
	return decode(uint64(score))
}

func ToScore(coord Coord) float64 {
	return float64(encode(coord))
}
