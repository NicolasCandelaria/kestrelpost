package resources

import (
	"errors"
	"strings"
)

type AntennaState string

const (
	AntennaTransmit AntennaState = "TRANSMIT"
	AntennaReceive  AntennaState = "RECEIVE"
	AntennaScan     AntennaState = "SCAN"
	AntennaCold     AntennaState = "COLD"
)

type Power struct {
	fuel              int
	capacity          int
	brownoutThreshold int
}

func NewPower(capacity, brownoutThreshold int) Power {
	if capacity <= 0 {
		capacity = 100
	}
	if brownoutThreshold < 0 || brownoutThreshold > capacity {
		brownoutThreshold = 20
	}
	return Power{fuel: capacity, capacity: capacity, brownoutThreshold: brownoutThreshold}
}

func (p *Power) Consume(units int) {
	if units <= 0 {
		return
	}
	p.fuel -= units
	if p.fuel < 0 {
		p.fuel = 0
	}
}

func (p *Power) Refuel(units int) {
	if units <= 0 {
		return
	}
	p.fuel += units
	if p.fuel > p.capacity {
		p.fuel = p.capacity
	}
}

func (p Power) Fuel() int {
	return p.fuel
}

func (p Power) Brownout() bool {
	return p.fuel <= p.brownoutThreshold
}

func (p Power) Gauge() string {
	switch {
	case p.fuel <= 0:
		return "[_____]"
	case p.fuel <= p.capacity/5:
		return "[#____]"
	case p.fuel <= (2*p.capacity)/5:
		return "[##___]"
	case p.fuel <= (3*p.capacity)/5:
		return "[###__]"
	case p.fuel <= (4*p.capacity)/5:
		return "[####_]"
	default:
		return "[#####]"
	}
}

type Antenna struct {
	state       AntennaState
	targetState AntennaState
	transition  bool
}

func NewAntenna() Antenna {
	return Antenna{state: AntennaReceive}
}

func (a *Antenna) BeginSwitch(target AntennaState) error {
	if a.transition {
		return errors.New("antenna already transitioning")
	}
	if target == a.state {
		return nil
	}
	a.targetState = target
	a.transition = true
	return nil
}

func (a *Antenna) CompleteSwitch() {
	if !a.transition {
		return
	}
	a.state = a.targetState
	a.transition = false
}

func (a Antenna) State() AntennaState {
	return a.state
}

func (a Antenna) Transitioning() bool {
	return a.transition
}

type Pins struct {
	max   int
	items []string
}

func NewPins(max int) Pins {
	if max <= 0 {
		max = 4
	}
	return Pins{max: max}
}

func (p *Pins) Pin(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" || p.IsPinned(id) {
		return true
	}
	if len(p.items) >= p.max {
		return false
	}
	p.items = append(p.items, id)
	return true
}

func (p *Pins) Unpin(id string) {
	idx := -1
	for i, item := range p.items {
		if item == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}
	p.items = append(p.items[:idx], p.items[idx+1:]...)
}

func (p Pins) IsPinned(id string) bool {
	for _, item := range p.items {
		if item == id {
			return true
		}
	}
	return false
}

func (p Pins) Items() []string {
	out := make([]string, len(p.items))
	copy(out, p.items)
	return out
}
