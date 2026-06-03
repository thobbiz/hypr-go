package layout

import "time"

type WindowEntry struct {
	Class     string `yaml:"class"`
	Workspace int    `yaml:"workspace"`
	X         int    `yaml:"x"`
	Y         int    `yaml:"y"`
	Width     int    `yaml:"width"`
	Height    int    `yaml:"height"`
	Floating  bool   `yaml:"floating"`
}

type Layout struct {
	Name    string        `yaml:"name"`
	SavedAt time.Time     `yaml:"saved_at"`
	Windows []WindowEntry `yaml:"windows"`
}
