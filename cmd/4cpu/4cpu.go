package main

import (
	"fmt"
	"time"

	"github.com/raffs/go-labs/sistema"
)

func header() {
	fmt.Println("CPU    User  |   Sys  |  io wait  |  guest | idle ")
}

func showCpu(nome string, c sistema.InfoCpuCore) {
	fmt.Printf("%s   %4.2f  |   %4.2f |   %4.2f   |  %4.2f  |   %4.2f\n", nome, c.User, c.System, c.IoWait, c.Guest, c.Idle)
}

func main() {
	sampleinfo1 := sistema.GetInfoCpuDefault()
	time.Sleep(500 * time.Millisecond)
	sampleinfo2 := sistema.GetInfoCpuDefault()

	header()
	for cpuNum := 0; cpuNum < len(sampleinfo2.Cores); cpuNum++ {
		info := sistema.InfoCpuCore{
			User:      sampleinfo2.Cores[cpuNum].User - sampleinfo1.Cores[cpuNum].User,
			Nice:      sampleinfo2.Cores[cpuNum].Nice - sampleinfo1.Cores[cpuNum].Nice,
			System:    sampleinfo2.Cores[cpuNum].System - sampleinfo1.Cores[cpuNum].System,
			Idle:      sampleinfo2.Cores[cpuNum].Idle - sampleinfo1.Cores[cpuNum].Idle,
			IoWait:    sampleinfo2.Cores[cpuNum].IoWait - sampleinfo1.Cores[cpuNum].IoWait,
			Irq:       sampleinfo2.Cores[cpuNum].Irq - sampleinfo1.Cores[cpuNum].Irq,
			SoftIrq:   sampleinfo2.Cores[cpuNum].SoftIrq - sampleinfo1.Cores[cpuNum].SoftIrq,
			Steal:     sampleinfo2.Cores[cpuNum].Steal - sampleinfo1.Cores[cpuNum].Steal,
			Guest:     sampleinfo2.Cores[cpuNum].Guest - sampleinfo1.Cores[cpuNum].Guest,
			GuestNice: sampleinfo2.Cores[cpuNum].GuestNice - sampleinfo1.Cores[cpuNum].GuestNice,
		}
		showCpu(fmt.Sprintf("cpu%d", cpuNum), info)
	}

	fmt.Println("----------")
	cpuUso := sistema.CpuUso()
	fmt.Printf("CPU Total: %.2f%%:\n", cpuUso)
}
