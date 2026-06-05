package geo

func New() *Geo {
	return &Geo{}
}

type Geo struct{}

func (g *Geo) Add(lon float64, lat float64, member string) bool {
	return true
}
