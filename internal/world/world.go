package world

type Weather string

const (
	WeatherClear    Weather = "CLEAR"
	WeatherOvercast Weather = "OVERCAST"
	WeatherStorm    Weather = "STORM"
)

type Fault string

const (
	FaultNone      Fault = "NONE"
	FaultIcing     Fault = "ANTENNA_ICING"
	FaultCapacitor Fault = "CAPACITOR"
	FaultFeedline  Fault = "FEEDLINE_FREEZE"
	FaultSurge     Fault = "GENERATOR_SURGE"
)

type Dog struct {
	Name   string
	Alive  bool
	Hunger int
}

func NewDog(name string) Dog {
	if name == "" {
		name = "Old Gray"
	}
	return Dog{Name: name, Alive: true}
}

func (d *Dog) Feed() {
	if !d.Alive {
		return
	}
	d.Hunger = 0
}

func (d *Dog) AdvanceNight() {
	if !d.Alive {
		return
	}
	d.Hunger++
	if d.Hunger >= 4 {
		d.Alive = false
	}
}

func WeatherForNight(night int) Weather {
	if night == 9 || night == 12 || night == 17 {
		return WeatherStorm
	}
	if night%3 == 0 {
		return WeatherOvercast
	}
	return WeatherClear
}

func FaultForNight(night int, w Weather) Fault {
	if night == 17 {
		return FaultSurge
	}
	if night == 9 && w == WeatherStorm {
		return FaultIcing
	}
	if w == WeatherStorm && night%2 == 0 {
		return FaultFeedline
	}
	if night%5 == 0 {
		return FaultCapacitor
	}
	return FaultNone
}
