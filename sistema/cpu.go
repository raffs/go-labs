package sistema

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/raffs/go-labs/fs"
)

const (
	CPU_STATFILE = "/proc/stat"

	// Campos do arquivo /proc/stat.
	//
	// Os valores sao computados com seu valor -1, pois,
	// os valores e extraido em uma `[]slice` separado
	// com todo os valores convertidos para `float64`.
	CPUSTAT_USER       = 0
	CPUSTAT_NICE       = 1
	CPUSTAT_SYSTEM     = 2
	CPUSTAT_IDLE       = 3
	CPUSTAT_IOWAIT     = 4
	CPUSTAT_IRQ        = 5
	CPUSTAT_SOFTIRQ    = 6
	CPUSTAT_STEAL      = 7
	CPUSTAT_GUEST      = 8
	CPUSTAT_GUEST_NICE = 9
)

var ClocksPerSec = float64(100)

/* Tempo que a CPU passou */
type InfoCpuCore struct {
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	IoWait    float64
	Irq       float64
	SoftIrq   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

type InfoCpu struct {
	Cpu            InfoCpuCore
	Cores          []InfoCpuCore
	Process        float64
	ProcessRunning float64
	ProccesBlocked float64
	SoftIRQ        float64
	ContextSwitch  float64
	BootTime       float64
}

func extrairNomeValores(linha string) (nome string, valores []float64) {
	colunas := strings.Fields(linha)

	for _, valorStr := range colunas[1:] {
		v, err := strconv.ParseFloat(valorStr, 64)
		if err == nil {
			valores = append(valores, v)
		} else {
			valores = append(valores, -1)
		}
	}

	nome = colunas[0]
	return nome, valores
}

func GetInfoCpu(caminhoArquivo string) (cpu InfoCpu) {
	linhas, err := fs.LerArquivo(caminhoArquivo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao abrir o arquivo '%s': %v\n", caminhoArquivo, err)
		return InfoCpu{}
	}

	for _, linha := range linhas {
		nome, colunas := extrairNomeValores(linha)
		if len(colunas) < 2 {
			break
		}

		if strings.HasPrefix(nome, "cpu") {
			info := InfoCpuCore{
				User:      colunas[CPUSTAT_USER] / ClocksPerSec,
				Nice:      colunas[CPUSTAT_NICE] / ClocksPerSec,
				System:    colunas[CPUSTAT_SYSTEM] / ClocksPerSec,
				Idle:      colunas[CPUSTAT_IDLE] / ClocksPerSec,
				IoWait:    colunas[CPUSTAT_IOWAIT] / ClocksPerSec,
				Irq:       colunas[CPUSTAT_IRQ] / ClocksPerSec,
				SoftIrq:   colunas[CPUSTAT_SOFTIRQ] / ClocksPerSec,
				Steal:     colunas[CPUSTAT_STEAL] / ClocksPerSec,
				Guest:     colunas[CPUSTAT_GUEST] / ClocksPerSec,
				GuestNice: colunas[CPUSTAT_GUEST_NICE] / ClocksPerSec,
			}
			if nome == "cpu" {
				cpu.Cpu = info
			} else {
				cpu.Cores = append(cpu.Cores, info)
			}
		} else if nome == "ctxt" {
			cpu.ContextSwitch = colunas[1]
		} else if nome == "btime" {
			cpu.BootTime = colunas[1]
		} else if nome == "processes" {
			cpu.Process = colunas[1]
		} else if nome == "procs_running" {
			cpu.ProcessRunning = colunas[1]
		} else if nome == "procs_blocked" {
			cpu.ProccesBlocked = colunas[1]
		} else if nome == "softirq" {
			cpu.SoftIRQ = colunas[1]
		}
	}

	return cpu
}

// GetInfoCpuDefault retorna uma número aleatorio, por enquanto.
func GetInfoCpuDefault() (info InfoCpu) {
	info = GetInfoCpu(CPU_STATFILE)
	return info
}

func CpuUso() float64 {
	cpuSampleStart := GetInfoCpuDefault()
	time.Sleep(500 * time.Millisecond)
	cpuSampleEnd := GetInfoCpuDefault()

	user := cpuSampleEnd.Cpu.User - cpuSampleStart.Cpu.User
	system := cpuSampleEnd.Cpu.System - cpuSampleStart.Cpu.System
	iowait := cpuSampleEnd.Cpu.IoWait - cpuSampleStart.Cpu.IoWait
	idle := cpuSampleEnd.Cpu.Idle - cpuSampleStart.Cpu.Idle

	active := user + system + iowait
	return active * 100 / (active + idle)
}
