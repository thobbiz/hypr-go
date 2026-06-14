package monitor

import "time"

type MonitorEntry struct {
	Name        string  `yaml:"name"`
	Width       int     `yaml:"width"`
	Height      int     `yaml:"height"`
	RefreshRate float64 `yaml:"refresh_rate"`
	PosX        int     `yaml:"pos_x"`
	PosY        int     `yaml:"pos_y"`
	Scale       float64 `yaml:"scale"`
	Transform   int     `yaml:"transform"`
}

type MonitorLayout struct {
	Name     string         `yaml:"name"`
	SavedAt  time.Time      `yaml:"saved_at"`
	Monitors []MonitorEntry `yaml:"monitors"`
}
