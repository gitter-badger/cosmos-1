package model

type Metric struct {
    Container string
    Cpu float32
}

type Metrics struct {
    Planet string
    Metrics []Metric
}
