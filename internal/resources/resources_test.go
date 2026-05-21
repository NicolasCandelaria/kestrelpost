package resources

import "testing"

func TestPinsCapAtFour(t *testing.T) {
	p := NewPins(4)
	if !p.Pin("maren") || !p.Pin("kid") || !p.Pin("harrow") || !p.Pin("cole") {
		t.Fatal("expected first four pins to succeed")
	}
	if p.Pin("osei") {
		t.Fatal("expected fifth pin to fail at cap")
	}
}

func TestAntennaTransition(t *testing.T) {
	a := NewAntenna()
	if err := a.BeginSwitch(AntennaScan); err != nil {
		t.Fatalf("begin switch: %v", err)
	}
	if !a.Transitioning() {
		t.Fatal("expected transition state")
	}
	a.CompleteSwitch()
	if a.State() != AntennaScan {
		t.Fatalf("state %s want %s", a.State(), AntennaScan)
	}
}

func TestPowerBrownout(t *testing.T) {
	p := NewPower(100, 20)
	p.Consume(85)
	if !p.Brownout() {
		t.Fatal("expected brownout")
	}
	if p.Gauge() != "[#____]" {
		t.Fatalf("gauge %q", p.Gauge())
	}
}
