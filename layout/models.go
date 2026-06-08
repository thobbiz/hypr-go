package layout

import "time"

type WindowEntry struct {
	Class     string `yaml:"class"`
	Cmd       string `yaml:"cmd"`
	Workspace int    `yaml:"workspace"`
	X         int    `yaml:"pos_x"`
	Y         int    `yaml:"pos_y"`
	Width     int    `yaml:"width"`
	Height    int    `yaml:"height"`
	Floating  bool   `yaml:"floating"`
}

type Layout struct {
	Name    string        `yaml:"name"`
	SavedAt time.Time     `yaml:"saved_at"`
	Windows []WindowEntry `yaml:"windows"`
}
