package world

import "testing"

func TestFixedFaultNights(t *testing.T) {
	if g := FaultForNight(9, WeatherStorm); g != FaultIcing {
		t.Fatalf("night 9 fault %s", g)
	}
	if g := FaultForNight(17, WeatherStorm); g != FaultSurge {
		t.Fatalf("night 17 fault %s", g)
	}
}

func TestDogDiesWhenUnfed(t *testing.T) {
	d := NewDog("Scout")
	for i := 0; i < 4; i++ {
		d.AdvanceNight()
	}
	if d.Alive {
		t.Fatal("expected dog to die after repeated hunger")
	}
}
