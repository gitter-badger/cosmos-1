package model

type MetricsContainer struct {
    Container string
    Cpu float32
}

type Metrics struct {
    Planet string
    Containers []*MetricsContainer
}
