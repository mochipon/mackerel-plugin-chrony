package chrony

import (
	"bytes"
	"flag"
	"log"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"chrony.stratum": {
		Label: "chrony Stratum",
		Unit: "integer",
		Metrics: []mp.Metrics{
			{Name: "stratum", Label: "Stratum", Diff: false, Stacked: false, Type: "integer"},
		},
	},
	"chrony.offset": {
		Label: "chrony Offset",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "last", Label: "Last (seconds)", Diff: false, Stacked: false, Type: "float64"},
			{Name: "rms", Label: "RMS (seconds)", Diff: false, Stacked: false, Type: "float64"},
		},
	},
	"chrony.system_time": {
		Label: "chrony System time",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "system_time", Label: "seconds", Diff: false, Stacked: false, Type: "float64"},
		},
	},

	"chrony.frequency": {
		Label: "chrony Frequency",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "frequency", Label: "ppm", Diff: false, Stacked: false, Type: "float64"},
		},
	},
	"chrony.residual_frequency": {
		Label: "chrony Residual freq",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "residual_frequency", Label: "ppm", Diff: false, Stacked: false, Type: "float64"},
		},
	},
	"chrony.skew": {
		Label: "chrony Skew",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "skew", Label: "ppm", Diff: false, Stacked: false, Type: "float64"},
		},
	},
	"chrony.root_delay": {
		Label: "chrony Root delay",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "root_delay", Label: "seconds", Diff: false, Stacked: false, Type: "float64"},
		},
	},
	"chrony.root_dispersion": {
		Label: "chrony Root dispersion",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "root_dispersion", Label: "seconds", Diff: false, Stacked: false, Type: "float64"},
		},
	},
	"chrony.update_interval": {
		Label: "chrony Update interval",
		Unit: "float",
		Metrics: []mp.Metrics{
			{Name: "update_interval", Label: "seconds", Diff: false, Stacked: false, Type: "float64"},
		},
	},
}

type chronyPlugin struct {
	path string
}

func (p chronyPlugin) fetchStats() (string, error) {
	cmd := exec.Command(p.path, "tracking")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (p chronyPlugin) FetchMetrics() (map[string]interface{}, error) {
	chronyStats := make(map[string]interface{})

	data, err := p.fetchStats()
	if err != nil {
		return nil, err
	}

	ignoreLine := []string{"ref", "status"}

	lines := strings.Split(strings.TrimSpace(data), "\n")
	for _, line := range lines {
		stats := strings.Split(line, ":")
		if len(stats) < 2 {
			log.Fatalf("unexpected output from chronyc, expected ':' in %s", data)
		}
		name := strings.ToLower(strings.Replace(strings.TrimSpace(stats[0]), " ", "_", -1))

		for _, v := range ignoreLine {
			if strings.Contains(name, v){
				continue
			}
		}

		valueFields := strings.Fields(stats[1])
		if len(valueFields) == 0 {
			log.Fatalf("unexpected output from chronyc: %s", data)
		}

		if strings.Contains(name, "stratum") {
			chronyStats["stratum"] = valueFields[0]
			continue
		}
		value, err := strconv.ParseFloat(valueFields[0], 64)
		if err != nil {
			continue
		}
		if strings.Contains(stats[1], "slow") {
			value = -value
		}

		var label string
		switch name {
			case "last_offset":
				label = "last"
			case "rms_offset":
				label = "rms"
			default:
				label = name
		}
		chronyStats[label] = value
	}

	return chronyStats, nil
}

func (p chronyPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func Do() {
	optCommand := flag.String("command", "/usr/bin/chronyc", "path to chronyc")

	flag.Parse()

	var chrony chronyPlugin

	chrony.path = *optCommand

	_, err := exec.LookPath(chrony.path)
	if err != nil {
		log.Fatalf("chronyc command is not found: %s", chrony.path)
	}

	helper := mp.NewMackerelPlugin(chrony)
	helper.Run()
}
