package monitor

import (
	"fmt"
	"log"
	"os"

	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func RestoreMonitor(name string) error {
	profile, err := decodeYaml(name)
	if err != nil {
		return err
	}
	for _, m := range profile.Monitors {
		if err := launchMonitor(m); err != nil {
			log.Printf("failed to restore monitor %s: %v", m.Name, err)
		}
	}
	return nil
}

func decodeYaml(filename string) (*MonitorLayout, error) {
	name := fmt.Sprintf("./data/monitor/%s.yaml", filename)
	f, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var layout MonitorLayout

	if err := yaml.Unmarshal(f, &layout); err != nil {
		return nil, err
	}

	return &layout, nil
}

func launchMonitor(m MonitorEntry) error {
	cmd := fmt.Sprintf("dispatch hl.dsp.exec_cmd('hyprctl keyword monitor %s,%dx%d@%.2f,%dx%d,%.2f')",
		m.Name, m.Width, m.Height, m.RefreshRate, m.PosX, m.PosY, m.Scale)
	resp, err := ipc.SendCommand(cmd)
	if err != nil {
		return err
	}
	if resp != "ok" {
		return fmt.Errorf("monitor restore failed for %s: %s", m.Name, resp)
	}

	return nil
}
