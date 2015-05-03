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

	return s
}

func (s *Series) AddPoint(point ...interface{}) {
	if s.Points == nil {
		s.Points = make([][]interface{}, 0)
	}
	s.Points = append(s.Points, point)
}
