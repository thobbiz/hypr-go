package monitor

import (
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func RestoreMonitor(name, dir string) error {
	profile, err := decodeYaml(name, dir)
	if err != nil {
		return err
	}

	var errs []error
	for _, m := range profile.Monitors {
		if err := launchMonitor(m); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func decodeYaml(name, dir string) (*MonitorLayout, error) {
	path, err := LayoutPath(dir, name)
	if err != nil {
		return nil, err
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, hyprerr.New("decodeYaml", hyprerr.ErrIO, err)
	}

	var layout MonitorLayout

	if err := yaml.Unmarshal(f, &layout); err != nil {
		return nil, hyprerr.New("decodeYaml", hyprerr.ErrParse, err)
	}

	return &layout, nil
}

func launchMonitor(m MonitorEntry) error {
	cmd := fmt.Sprintf("dispatch hl.dsp.exec_cmd('hyprctl keyword monitor %s,%dx%d@%.2f,%dx%d,%.2f,transform,%d')",
		m.Name, m.Width, m.Height, m.RefreshRate, m.PosX, m.PosY, m.Scale, m.Transform)

	resp, err := ipc.SendCommand(cmd)
	if err != nil {
		return hyprerr.New("launchMonitor", hyprerr.ErrIPC, fmt.Errorf("monitor %s: %w", m.Name, err))
	}
	if resp != "ok" {
		return hyprerr.New("launchMonitor", hyprerr.ErrIPC, fmt.Errorf("monitor %s: dispatch rejected: %s", m.Name, resp))
	}

	return confirmMonitor(m)
}

func confirmMonitor(want MonitorEntry) error {
	for range 10 {
		time.Sleep(100 * time.Millisecond)

		current, err := getMonitors()
		if err != nil {
			return err
		}
		for _, got := range *current {
			if got.Name == want.Name && sameState(want, got) {
				return nil
			}
		}
	}
	return hyprerr.New("confirmMonitor", hyprerr.ErrIPC,
		fmt.Errorf("monitor %s did not reach the saved state", want.Name))
}

func sameState(a, b MonitorEntry) bool {
	return a.Width == b.Width && a.Height == b.Height &&
		a.PosX == b.PosX && a.PosY == b.PosY &&
		a.Transform == b.Transform &&
		math.Abs(a.Scale-b.Scale) < 0.01 &&
		math.Abs(a.RefreshRate-b.RefreshRate) < 0.05
}
