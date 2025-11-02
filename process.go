package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
)

var (
	processLines = make(map[string]int)
	lineMutex    sync.Mutex
	nextLine     = 0
	initialized  = false
)

type Process struct {
	uniqueID string
	name     string
	state    string
	percent  int
	finish   bool
}

func NewProcess() *Process {
	return &Process{
		uniqueID: uuid.New().String(),
	}
}

func (p *Process) Name() string {
	return p.name
}

func (p *Process) SetName(name string) {
	p.name = name
}

func (p *Process) Finish() {
	p.finish = true
}

func (p *Process) Load(state string, percent int) *Process {
	p.state = state
	p.percent = percent
	return p
}

func (p *Process) Process() {
	lineMutex.Lock()
	defer lineMutex.Unlock()

	if !initialized {
		fmt.Print("\033[2J\033[H")
		initialized = true
	}

	line, exists := processLines[p.uniqueID]
	if !exists {
		line = nextLine
		processLines[p.uniqueID] = line
		nextLine++

		if line > 0 {
			fmt.Printf("\033[%dB", line)
		}
	}

	var load int
	var remained int

	switch {
	case p.percent == 0:
		load = 0
		remained = 30
	case p.percent == 100:
		load = 30
		remained = 0
	case p.percent > 0 && p.percent < 100:
		load = int(float32(p.percent) * 30.00 / 100)
		remained = 30 - load
	}

	progress := fmt.Sprintf("%s[%d%% %s>%s ~ %s]",
		p.name,
		p.percent,
		strings.Repeat("=", load),
		strings.Repeat(".", remained),
		p.state,
	)

	fmt.Printf("\033[%d;0H", line+1)
	fmt.Printf("\r%-80s", progress)

	fmt.Printf("\033[%d;0H", nextLine+1)

	if p.finish {
		delete(processLines, p.uniqueID)
	}
}

func (p *Process) Cleanup() {
	lineMutex.Lock()
	defer lineMutex.Unlock()

	if line, exists := processLines[p.uniqueID]; exists {
		delete(processLines, p.uniqueID)

		if p.finish {
			fmt.Printf("\033[%d;0H", line+1)
			fmt.Printf("\r%-80s", "")
		}

		fmt.Printf("\033[%d;0H", nextLine+1)
	}
}
