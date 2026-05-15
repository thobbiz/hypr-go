package main

type Workspace struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Monitor      string `json:"monitor"`
	MonitorID    int    `json:"monitorID"`
	Windows      int    `json:"windows"` // Number of open windows
	IsPersistent bool   `json:"ispersistent"`
	TiledLayout  string `json:"tiledLayout"` // Can be dwindle or master
}
