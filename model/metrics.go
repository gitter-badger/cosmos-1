package model

type MetricsContainerParam struct {
    Container string
    Cpu float32
    Memory uint64
}

type MetricsParam struct {
    Planet string
    Containers []*MetricsContainerParam
}
