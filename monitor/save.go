package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func SaveMonitor(name string) error {
	monitor, err := getMonitor(name)
	if err != nil {
		return err
	}

	monitorLayout := MonitorLayout{
		Name:     name,
		SavedAt:  time.Now(),
		Monitors: *monitor,
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
		log.Fatal(err)
	}

	name := fmt.Sprintf("./data/monitor/%s.yaml", monitorLayout.Name)
	err = os.WriteFile(name, data, 0644)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}

	return nil
}

func getMonitor(name string) (*[]MonitorEntry, error) {
	resp, err := ipc.SendCommand("j/monitors")
	if err != nil {
		log.Fatal(err)
	}

	var monitors []ipc.Monitor
	if err := json.Unmarshal([]byte(resp), &monitors); err != nil {
		return nil, err
	}

	var result []MonitorEntry

	for _, m := range monitors {
		if m.Name == name {
			anon := MonitorEntry{
				Name:        name,
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
	}

	if len(result) == 0 {
		return nil, errors.New("monitor not found")
	}
	return &result, nil
}
