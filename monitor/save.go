package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func SaveMonitor(name string) error {
	monitors, err := getMonitors()
	if err != nil {
		return err
	}

	monitorLayout := MonitorLayout{
		Name:     name,
		SavedAt:  time.Now(),
		Monitors: *monitors,
	}

	err = saveMonitorAsYaml(monitorLayout)
	if err != nil {
		return err
	}

	return nil
}

func saveMonitorAsYaml(monitorLayout MonitorLayout) error {
	data, err := yaml.Marshal(monitorLayout)
	if err != nil {
		return hyprerr.New("saveMonitorAsYaml", hyprerr.ErrParse, err)
	}

	name := fmt.Sprintf("./data/monitor/%s.yaml", monitorLayout.Name)
	err = os.WriteFile(name, data, 0644)
	if err != nil {
		return hyprerr.New("saveMonitorAsYaml", hyprerr.ErrIO, err)
	}

	return nil
}

func getMonitors() (*[]MonitorEntry, error) {
	resp, err := ipc.SendCommand("j/monitors")
	if err != nil {
		return nil, hyprerr.New("getMonitors", hyprerr.ErrIPC, err)
	}

	var monitors []ipc.Monitor
	if err := json.Unmarshal([]byte(resp), &monitors); err != nil {
		return nil, hyprerr.New("getMonitors", hyprerr.ErrParse, err)
	}

	var result []MonitorEntry

	for _, m := range monitors {
		anon := MonitorEntry{
			Name:        m.Name,
			Width:       m.Width,
			Height:      m.Height,
			RefreshRate: m.RefreshRate,
			PosX:        m.X,
			PosY:        m.Y,
			Scale:       m.Scale,
			Transform:   m.Transform,
		}
		result = append(result, anon)
	}

	if len(result) == 0 {
		return nil, hyprerr.New("getMonitors", hyprerr.ErrNotFound, fmt.Errorf("no monitors found"))
	}
	return &result, nil
}
