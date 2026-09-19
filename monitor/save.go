package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func SaveMonitor(name string, dir string) error {
	path, err := LayoutPath(dir, name)
	if err != nil {
		return err
	}

	monitors, err := getMonitors()
	if err != nil {
		return err
	}

	monitorLayout := MonitorLayout{
		Name:     name,
		SavedAt:  time.Now(),
		Monitors: *monitors,
	}

	err = saveMonitorLayoutAsYaml(monitorLayout, path)
	if err != nil {
		return err
	}

	return nil
}

func saveMonitorLayoutAsYaml(layout MonitorLayout, path string) error {
	data, err := yaml.Marshal(layout)
	if err != nil {
		return hyprerr.New("saveMonitorLayoutAsYaml", hyprerr.ErrParse, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return hyprerr.New("saveMonitorLayoutAsYaml", hyprerr.ErrIO, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return hyprerr.New("saveMonitorLayoutAsYaml", hyprerr.ErrIO, err)
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

func LayoutPath(dir, name string) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}

	base, err := resolveDir(dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(base, name+".yaml"), nil
}

func validateName(name string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) || strings.HasSuffix(name, ".yaml") {
		return hyprerr.New("validateName", hyprerr.ErrInvalidInput, fmt.Errorf("invalid layout name %q", name))
	}
	return nil
}

func resolveDir(dir string) (string, error) {
	if dir == "" {
		return defaultDir()
	}

	if dir == "~" || strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", hyprerr.New("resolveDir", hyprerr.ErrIO, err)
		}
		dir = filepath.Join(home, strings.TrimPrefix(dir, "~"))
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", hyprerr.New("resolveDir", hyprerr.ErrIO, err)
	}
	return abs, nil
}

func defaultDir() (string, error) {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", hyprerr.New("defaultDir", hyprerr.ErrIO, err)
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "hypr-go", "monitor"), nil
}
