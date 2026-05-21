package content

import "testing"

func TestLoadCampaignTwentyNights(t *testing.T) {
	c, err := LoadCampaign()
	if err != nil {
		t.Fatalf("load campaign: %v", err)
	}
	if len(c.Nights) != 20 {
		t.Fatalf("nights %d want 20", len(c.Nights))
	}
	if c.Nights[1].Source != "MAREN" {
		t.Fatalf("night 1 source %s", c.Nights[1].Source)
	}
}
