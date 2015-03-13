package influxdbc

type Series struct {
	Name    string          `json:"name"`
	Columns []string        `json:"columns"`
	Points  [][]interface{} `json:"points"`
}

func NewSeries(name string, cols ...string) *Series {
	s := new(Series)
	s.Name = name
	s.Columns = cols
	s.Points = make([][]interface{}, 0)
	return s
}

func (s *Series) AddPoint(point ...interface{}) {
	s.Points = append(s.Points, point)
}
