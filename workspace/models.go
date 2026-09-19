package workspace

import "time"

type WindowEntry struct {
	Class     string `yaml:"class"`
	Exec      string `yaml:"exec"`
	Workspace int    `yaml:"workspace"`
	X         int    `yaml:"pos_x"`
	Y         int    `yaml:"pos_y"`
	Width     int    `yaml:"width"`
	Height    int    `yaml:"height"`
	Floating  bool   `yaml:"floating"`
}

type WorkspaceLayout struct {
	Name    string        `yaml:"name"`
	SavedAt time.Time     `yaml:"saved_at"`
	Windows []WindowEntry `yaml:"windows"`
}
