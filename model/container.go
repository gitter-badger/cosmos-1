package model

import "fmt"

type Port struct {
	PrivatePort *int
	PublicPort  *int
	Type        *string
}

func (p *Port) Description() string {
	return fmt.Sprintf("%d:%d %s", *p.PublicPort, *p.PrivatePort, *p.Type)
}

type Network struct {
	RxBytes *int64
	TxBytes *int64
}

type Cpu struct {
	TotalUtilization  *float32
	PerCpuUtilization []float32
}

type Memory struct {
	Limit *int64
	Usage *int64
}

type Stats struct {
	Network *Network
	Cpu     *Cpu
	Memory  *Memory
}

type Container struct {
	Id         *string
	Command    *string
	Image      *string
	Names      []string
	Ports      []*Port
	Status     *string
	SizeRw     *int64
	SizeRootFs *int64
	Stats      *Stats
}
