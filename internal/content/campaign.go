package content

import (
	"embed"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

//go:embed data/*.yaml
var campaignFS embed.FS

type Effects struct {
	Fuel         int  `yaml:"fuel"`
	Hub          int  `yaml:"hub"`
	Trust        int  `yaml:"trust"`
	Kid          int  `yaml:"kid"`
	SetHarrow    bool `yaml:"set_harrow"`
	OseiRelease  bool `yaml:"osei_release"`
	ConvoyBetray bool `yaml:"convoy_betray"`
}

type Choice struct {
	ID      string  `yaml:"id"`
	Text    string  `yaml:"text"`
	Effects Effects `yaml:"effects"`
}

type Night struct {
	Number   int      `yaml:"number"`
	Act      int      `yaml:"act"`
	Source   string   `yaml:"source"`
	Hash     string   `yaml:"hash"`
	PreShift string   `yaml:"pre_shift"`
	Quote    string   `yaml:"quote"`
	Incident string   `yaml:"incident"`
	Choices  []Choice `yaml:"choices"`
}

type Campaign struct {
	Nights map[int]Night
	Order  []int
}

func LoadCampaign() (Campaign, error) {
	entries, err := campaignFS.ReadDir("data")
	if err != nil {
		return Campaign{}, err
	}

	c := Campaign{Nights: make(map[int]Night)}
	for _, entry := range entries {
		raw, readErr := campaignFS.ReadFile("data/" + entry.Name())
		if readErr != nil {
			return Campaign{}, readErr
		}
		var n Night
		if unmarshalErr := yaml.Unmarshal(raw, &n); unmarshalErr != nil {
			return Campaign{}, fmt.Errorf("%s: %w", entry.Name(), unmarshalErr)
		}
		if n.Number == 0 {
			return Campaign{}, fmt.Errorf("%s: missing number", entry.Name())
		}
		if len(n.Choices) != 3 {
			return Campaign{}, fmt.Errorf("%s: expected 3 choices", entry.Name())
		}
		c.Nights[n.Number] = n
		c.Order = append(c.Order, n.Number)
	}
	sort.Ints(c.Order)
	return c, nil
}
